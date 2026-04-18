package api

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"

	"github.com/rampartfw/rampart/internal/api/handlers"
	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/cluster"
	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/snapshot"
)

//go:embed ui/dist/*
var uiFS embed.FS

type contextKey string

const paramsKey = handlers.ParamsKey

type Middleware func(http.Handler) http.Handler

type Router struct {
	trees       map[string]*node
	middlewares []Middleware
}

type node struct {
	path     string
	handler  http.Handler
	children []*node
	param    string
	isParam  bool
}

func NewRouter() *Router {
	return &Router{
		trees: make(map[string]*node),
	}
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) Handle(method, pattern string, handler http.HandlerFunc) {
	h := http.Handler(handler)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}
	r.addRoute(method, pattern, h)
}

func (r *Router) addRoute(method, pattern string, handler http.Handler) {
	root, ok := r.trees[method]
	if !ok {
		root = &node{path: "/"}
		r.trees[method] = root
	}

	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	if pattern == "/" {
		parts = []string{""}
	}

	current := root
	for _, part := range parts {
		if part == "" && pattern != "/" {
			continue
		}

		var found *node
		for _, child := range current.children {
			if child.path == part {
				found = child
				break
			}
		}

		if found == nil {
			newNode := &node{path: part}
			if strings.HasPrefix(part, ":") {
				newNode.isParam = true
				newNode.param = part[1:]
			}
			current.children = append(current.children, newNode)
			found = newNode
		}
		current = found
	}
	current.handler = handler
}

func (r *Router) find(method, path string) (http.Handler, map[string]string) {
	root, ok := r.trees[method]
	if !ok {
		return nil, nil
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if path == "/" {
		parts = []string{""}
	}

	params := make(map[string]string)
	current := root

Outer:
	for _, part := range parts {
		if part == "" && path != "/" {
			continue
		}

		for _, child := range current.children {
			if child.path == part || child.isParam {
				if child.isParam {
					params[child.param] = part
				}
				current = child
				continue Outer
			}
		}
		return nil, nil
	}

	return current.handler, params
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, params := r.find(req.Method, req.URL.Path)
	if handler == nil {
		http.Error(w, "Not Found", 404)
		return
	}

	ctx := context.WithValue(req.Context(), paramsKey, params)
	handler.ServeHTTP(w, req.WithContext(ctx))
}

func Params(r *http.Request) map[string]string {
	if params, ok := r.Context().Value(paramsKey).(map[string]string); ok {
		return params
	}
	return nil
}

type Server struct {
	cfg           *config.Config
	engine        *engine.Engine
	snapshotStore *snapshot.Store
	auditStore    *audit.Store
	raftNode      *cluster.RaftNode
	router        *Router
}

func NewServer(cfg *config.Config, eng *engine.Engine, ss *snapshot.Store, as *audit.Store, rn *cluster.RaftNode) *Server {
	s := &Server{
		cfg:           cfg,
		engine:        eng,
		snapshotStore: ss,
		auditStore:    as,
		raftNode:      rn,
		router:        NewRouter(),
	}

	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) serveUI() http.Handler {
	// ui/dist contains the build output
	distFS, err := fs.Sub(uiFS, "ui/dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the path relative to /ui/
		path := r.URL.Path
		if !strings.HasPrefix(path, "/ui/") && path != "/ui" {
			http.NotFound(w, r)
			return
		}

		// 2. Clean the path to get the relative file path
		relPath := strings.TrimPrefix(path, "/ui/")
		if relPath == "/ui" || relPath == "" {
			relPath = "index.html"
		}

		// 3. Check if the requested file exists in the embedded FS
		f, err := distFS.Open(relPath)
		if err != nil {
			// If file doesn't exist, it might be a React Router path - serve index.html
			r.URL.Path = "/" // For http.FileServer to serve index.html from root of distFS
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()

		// 4. Serve the actual file
		// We need to set r.URL.Path to relPath so fileServer finds it at the root of distFS
		r.URL.Path = "/" + relPath
		fileServer.ServeHTTP(w, r)
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "error",
		"error": map[string]string{
			"message": message,
		},
	})
}

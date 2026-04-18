package api

import (
	"crypto/subtle"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/model"
	"github.com/rampartfw/rampart/internal/security"
)

var apiLimiter = security.NewRateLimiter(10, 20) // 10 req/sec, burst 20

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC recovered in API: %v\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = model.GenerateUUIDv7()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(keys []config.APIKey) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for unix socket (simplified for now)
			if strings.HasPrefix(r.RemoteAddr, "@") || strings.Contains(r.RemoteAddr, "/") {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			valid := false
			var matchedKey string
			for _, k := range keys {
				if subtle.ConstantTimeCompare([]byte(k.Key), []byte(token)) == 1 {
					valid = true
					matchedKey = k.Name
					break
				}
			}

			if !valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Rate Limiting
			if !apiLimiter.Allow(matchedKey) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

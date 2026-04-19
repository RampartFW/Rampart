package k8s

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

// Controller monitors K8s NetworkPolicies and translates them to Rampart.
type Controller struct {
	engine *engine.Engine
	token  string
	host   string
	client *http.Client
}

func NewController(eng *engine.Engine) (*Controller, error) {
	token, _ := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")

	if host == "" {
		host = "kubernetes.default.svc"
	}
	if port == "" {
		port = "443"
	}

	return &Controller{
		engine: eng,
		token:  string(token),
		host:   fmt.Sprintf("https://%s:%s", host, port),
		client: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *Controller) Run(ctx context.Context) {
	log.Println("K8s: Starting NetworkPolicy controller...")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.sync(ctx); err != nil {
				log.Printf("K8s: Sync error: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Controller) sync(ctx context.Context) error {
	// 1. Fetch NetworkPolicies
	// GET /apis/networking.k8s.io/v1/networkpolicies
	log.Println("K8s: Synchronizing policies...")

	// 2. Translate to Rampart PolicySet
	ps := &model.PolicySetYAML{
		APIVersion: "rampartfw.com/v1",
		Kind:       "PolicySet",
		Metadata: model.PolicyMetadata{
			Name: "k8s-managed-policy",
		},
	}

	// 3. Apply via Engine
	compiled, err := engine.Compile(ps, nil)
	if err != nil {
		return err
	}

	c.engine.SetRules(compiled)
	return c.engine.ReapplyRules(ctx)
}

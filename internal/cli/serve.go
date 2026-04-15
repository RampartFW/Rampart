package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rampartfw/rampart/internal/api"
	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/logger"
	"github.com/rampartfw/rampart/internal/snapshot"
)

type ServeCommand struct{}

func (c *ServeCommand) Name() string        { return "serve" }
func (c *ServeCommand) Description() string { return "Start server (API + WebUI + Raft)" }

func (c *ServeCommand) Run(args []string) {
	// 1. Load configuration
	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	if err := logger.Init(cfg); err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	log := logger.Component("cli")

	log.Info("Starting Rampart server", "version", "1.0.0")

	// 3. Initialize backend
	var b backend.Backend
	
	// Convert config.BackendConfig to backend.BackendConfig
	bcfg := backend.BackendConfig{
		Type:     cfg.Backend.Type,
		Settings: make(map[string]string),
	}
	
	// Fill settings based on backend type
	switch cfg.Backend.Type {
	case "nftables":
		bcfg.Settings["tableName"] = cfg.Backend.Nftables.TableName
		bcfg.Settings["binary"] = cfg.Backend.Nftables.Binary
	case "iptables":
		bcfg.Settings["chainPrefix"] = cfg.Backend.Iptables.ChainPrefix
		bcfg.Settings["binary"] = cfg.Backend.Iptables.Binary
	case "aws":
		bcfg.Settings["region"] = cfg.Backend.AWS.Region
		bcfg.Settings["securityGroupId"] = cfg.Backend.AWS.SecurityGroupId
	}

	if cfg.Backend.Type == "auto" {
		b, err = backend.AutoDetect()
	} else {
		b, err = backend.NewBackend(cfg.Backend.Type, bcfg)
	}
	if err != nil {
		log.Error("Failed to initialize backend", "error", err)
		os.Exit(1)
	}
	log.Info("Backend initialized", "type", b.Name())

	// 4. Initialize stores
	auditStore, err := audit.NewStore(cfg.Audit.Directory, cfg.Audit.Retention)
	if err != nil {
		log.Error("Failed to initialize audit store", "error", err)
		os.Exit(1)
	}

	snapshotStore, err := snapshot.NewStore(cfg.Snapshots.Directory)
	if err != nil {
		log.Error("Failed to initialize snapshot store", "error", err)
		os.Exit(1)
	}

	// 5. Initialize engine
	eng := engine.NewEngine(b)
	
	// Load current state from backend
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	current, err := b.CurrentState(ctx)
	if err != nil {
		log.Warn("Failed to load current state from backend", "error", err)
	} else {
		eng.SetRules(current)
		log.Info("Current state loaded from backend", "rules", len(current.Rules))
	}

	// 6. Initialize API server
	srv := api.NewServer(cfg, eng, snapshotStore, auditStore)

	httpServer := &http.Server{
		Addr:    cfg.Server.Listen,
		Handler: srv,
	}

	// 7. Start server in a goroutine
	go func() {
		log.Info("API server listening", "addr", cfg.Server.Listen)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 8. Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Info("Shutting down server...")

	// 9. Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped")
}

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
	{
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		current, err := b.CurrentState(ctx)
		cancel()
		if err != nil {
			log.Warn("Failed to load current state from backend", "error", err)
		} else {
			eng.SetRules(current)
			log.Info("Current state loaded from backend", "rules", len(current.Rules))
		}
	}
	
	// Start scheduler for time-based rules
	sched := engine.NewScheduler(eng, cfg.Scheduler.CheckInterval)
	go sched.Run(context.Background())
	log.Info("Scheduler started", "interval", cfg.Scheduler.CheckInterval)

	// 6. Initialize Cluster (Raft)
	var raftNode *cluster.RaftNode
	if cfg.Cluster.Enabled {
		log.Info("Initializing cluster node", "id", cfg.Cluster.NodeID)
		
		// Setup FSM
		fsm := cluster.NewPolicyFSM(b, eng)
		
		// Setup Log
		logDir := filepath.Join(cfg.Audit.Directory, "raft")
		os.MkdirAll(logDir, 0755)
		raftLog, err := cluster.NewLog(filepath.Join(logDir, "wal.gob"))
		if err != nil {
			log.Error("Failed to initialize Raft log", "error", err)
			os.Exit(1)
		}
		
		// Setup Transport
		transport, err := cluster.NewTCPTransport(
			cfg.Cluster.TLS.Cert,
			cfg.Cluster.TLS.Key,
			cfg.Cluster.TLS.CA,
		)
		if err != nil {
			log.Error("Failed to initialize Raft transport", "error", err)
			os.Exit(1)
		}
		
		// Setup Peers
		peers := make(map[string]string)
		for i, peerAddr := range cfg.Cluster.Peers {
			// Simplified: using index as part of ID if not specified
			// In a real prod environment, IDs should be stable and unique
			peerID := fmt.Sprintf("node-%d", i+1)
			peers[peerID] = peerAddr
		}
		// Ensure self is in peers or handled by RaftNode
		
		raftNode = cluster.NewRaftNode(cfg.Cluster.NodeID, peers, transport, raftLog, fsm)
		
		if err := raftNode.Start(cfg.Cluster.Listen); err != nil {
			log.Error("Failed to start Raft node", "error", err)
			os.Exit(1)
		}
		log.Info("Cluster node started", "listen", cfg.Cluster.Listen, "advertise", cfg.Cluster.Advertise)
	}

	// 7. Initialize API server
	srv := api.NewServer(cfg, eng, snapshotStore, auditStore, raftNode)

	httpServer := &http.Server{
		Addr:    cfg.Server.Listen,
		Handler: srv,
	}

	// 8. Start server in a goroutine (HTTPS)
	go func() {
		log.Info("API server listening (HTTPS)", "addr", cfg.Server.Listen)
		// We use the same TLS certs as the cluster for simplicity and consistency
		if err := httpServer.ListenAndServeTLS(cfg.Cluster.TLS.Cert, cfg.Cluster.TLS.Key); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 9. Wait for termination or reload signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	for {
		sig := <-stop
		if sig == syscall.SIGHUP {
			log.Info("Reloading configuration...")
			newCfg, err := config.LoadConfig(ConfigPath)
			if err != nil {
				log.Error("Failed to reload configuration", "error", err)
				continue
			}
			
			// Update logger level/format dynamically
			logger.Init(newCfg)
			
			// Note: Some changes like server listen address or cluster node ID 
			// still require restart. We log a warning for those.
			if newCfg.Server.Listen != cfg.Server.Listen {
				log.Warn("Server listen address change requires restart to take effect")
			}
			
			cfg = newCfg
			log.Info("Configuration reloaded successfully")
			continue
		}
		
		log.Info("Shutting down server...", "signal", sig)
		break
	}

	// 10. Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if raftNode != nil {
		raftNode.Close()
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped")
}

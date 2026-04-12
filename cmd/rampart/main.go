package main

import (
	"log"

	"github.com/rampartfw/rampart/internal/cli"
	"github.com/rampartfw/rampart/internal/security"

	_ "github.com/rampartfw/rampart/internal/backend/aws"
	_ "github.com/rampartfw/rampart/internal/backend/azure"
	_ "github.com/rampartfw/rampart/internal/backend/ebpf"
	_ "github.com/rampartfw/rampart/internal/backend/gcp"
	_ "github.com/rampartfw/rampart/internal/backend/iptables"
	_ "github.com/rampartfw/rampart/internal/backend/nftables"
)

func main() {
	// Initialize security (drop capabilities if possible)
	if err := security.Initialize(); err != nil {
		log.Printf("Warning: security initialization failed: %v", err)
	}

	cli.Execute()
}

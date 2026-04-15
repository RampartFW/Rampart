package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/snapshot"
)

type SnapshotCommand struct{}

func (c *SnapshotCommand) Name() string        { return "snapshot" }
func (c *SnapshotCommand) Description() string { return "Manage configuration snapshots" }

func (c *SnapshotCommand) Run(args []string) {
	if len(args) == 0 {
		c.usage()
		return
	}

	sub := args[0]
	subArgs := args[1:]

	switch sub {
	case "list":
		c.list(subArgs)
	case "create":
		c.create(subArgs)
	case "delete", "remove":
		c.delete(subArgs)
	case "help", "-h", "--help":
		c.usage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", sub)
		c.usage()
		os.Exit(1)
	}
}

func (c *SnapshotCommand) usage() {
	fmt.Println("Usage: rampart snapshot <subcommand> [arguments]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  list      List all available snapshots")
	fmt.Println("  create    Create a new snapshot from current state")
	fmt.Println("  delete    Remove a specific snapshot")
}

func (c *SnapshotCommand) getBackend(cfg *config.Config) (backend.Backend, error) {
	if cfg.Backend.Type == "auto" {
		return backend.AutoDetect()
	}

	bcfg := backend.BackendConfig{
		Type:     cfg.Backend.Type,
		Settings: make(map[string]string),
	}
	switch cfg.Backend.Type {
	case "nftables":
		bcfg.Settings["tableName"] = cfg.Backend.Nftables.TableName
	case "iptables":
		bcfg.Settings["chainPrefix"] = cfg.Backend.Iptables.ChainPrefix
	case "aws":
		bcfg.Settings["region"] = cfg.Backend.AWS.Region
		bcfg.Settings["securityGroupId"] = cfg.Backend.AWS.SecurityGroupId
	case "mock":
		// No special settings
	}
	return backend.NewBackend(cfg.Backend.Type, bcfg)
}

func (c *SnapshotCommand) list(args []string) {
	fs := flag.NewFlagSet("snapshot list", flag.ExitOnError)
	fs.Parse(args)

	cfg := LoadConfig()
	snapStore, err := snapshot.NewStore(cfg.Snapshots.Directory)
	ExitOnError(err, "Initialize snapshot store")
	
	snaps, err := snapStore.List()
	ExitOnError(err, "List snapshots")

	if len(snaps) == 0 {
		fmt.Println("No snapshots found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTRIGGER\tCREATED\tDESCRIPTION")
	for _, s := range snaps {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
			s.ID[:8], s.Trigger, s.CreatedAt.Format("2006-01-02 15:04"), s.Description)
	}
	w.Flush()
}

func (c *SnapshotCommand) create(args []string) {
	fs := flag.NewFlagSet("snapshot create", flag.ExitOnError)
	desc := fs.String("desc", "Manual snapshot", "Description")
	fs.Parse(args)

	cfg := LoadConfig()
	be, err := c.getBackend(cfg)
	ExitOnError(err, "Initialize backend")
	defer be.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	current, err := be.CurrentState(ctx)
	ExitOnError(err, "Get current state")

	snapStore, err := snapshot.NewStore(cfg.Snapshots.Directory)
	ExitOnError(err, "Initialize snapshot store")
	
	s, err := snapStore.Create("manual", *desc, current)
	ExitOnError(err, "Create snapshot")

	fmt.Printf("%s Created snapshot %s\n", Colorize("✓", colorGreen), s.ID)
}

func (c *SnapshotCommand) delete(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Snapshot ID is required")
		os.Exit(1)
	}
	id := args[0]

	cfg := LoadConfig()
	snapStore, err := snapshot.NewStore(cfg.Snapshots.Directory)
	ExitOnError(err, "Initialize snapshot store")
	
	err = snapStore.Delete(id)
	ExitOnError(err, "Delete snapshot")

	fmt.Printf("%s Deleted snapshot %s\n", Colorize("✓", colorGreen), id)
}

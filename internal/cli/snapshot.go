package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/snapshot"
)

type SnapshotCommand struct{}

func (c *SnapshotCommand) Name() string        { return "snapshot" }
func (c *SnapshotCommand) Description() string { return "Manage snapshots" }

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
	case "delete":
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
	fmt.Println("  list      List all snapshots")
	fmt.Println("  create    Create a new snapshot")
	fmt.Println("  delete    Delete a snapshot")
}

func (c *SnapshotCommand) list(args []string) {
	fs := flag.NewFlagSet("snapshot list", flag.ExitOnError)
	snapDir := fs.String("dir", "snapshots", "Snapshots directory")
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	snapStore, err := snapshot.NewStore(*snapDir)
	ExitOnError(err, "Initialize snapshot store")
	snaps, err := snapStore.List()
	ExitOnError(err, "List snapshots")

	fmt.Printf("%-24s %-20s %s\n", "ID", "Created", "Description")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, s := range snaps {
		fmt.Printf("%-24s %-20s %s\n", s.ID, s.CreatedAt.Format("2006-01-02 15:04:05"), s.Description)
	}
}

func (c *SnapshotCommand) create(args []string) {
	fs := flag.NewFlagSet("snapshot create", flag.ExitOnError)
	snapDir := fs.String("dir", "snapshots", "Snapshots directory")
	desc := fs.String("desc", "Manual snapshot", "Description")
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	be, err := backend.AutoDetect()
	ExitOnError(err, "Auto-detect backend")

	current, err := be.CurrentState()
	ExitOnError(err, "Get current state")

	if _, err := os.Stat(*snapDir); os.IsNotExist(err) {
		os.MkdirAll(*snapDir, 0755)
	}
	snapStore, err := snapshot.NewStore(*snapDir)
	ExitOnError(err, "Initialize snapshot store")
	s, err := snapStore.Create("manual", *desc, current)
	ExitOnError(err, "Create snapshot")

	fmt.Printf("%s Created snapshot %s\n", Colorize("✓", colorGreen), s.ID)
}

func (c *SnapshotCommand) delete(args []string) {
	fs := flag.NewFlagSet("snapshot delete", flag.ExitOnError)
	snapDir := fs.String("dir", "snapshots", "Snapshots directory")
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	if len(fs.Args()) == 0 {
		fmt.Fprintln(os.Stderr, Colorize("Error", colorRed)+": Snapshot ID is required")
		os.Exit(1)
	}
	id := fs.Args()[0]

	snapStore, err := snapshot.NewStore(*snapDir)
	ExitOnError(err, "Initialize snapshot store")
	err = snapStore.Delete(id)
	ExitOnError(err, "Delete snapshot")

	fmt.Printf("%s Deleted snapshot %s\n", Colorize("✓", colorGreen), id)
}

package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/snapshot"
)

type RollbackCommand struct{}

func (c *RollbackCommand) Name() string        { return "rollback" }
func (c *RollbackCommand) Description() string { return "Rollback to a snapshot" }

func (c *RollbackCommand) Run(args []string) {
	fs := flag.NewFlagSet("rollback", flag.ExitOnError)
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
	s, rs, err := snapStore.Load(id)
	ExitOnError(err, "Load snapshot")

	be, err := backend.AutoDetect()
	ExitOnError(err, "Auto-detect backend")

	fmt.Printf("Rolling back to snapshot %s (%s)...\n", s.ID, s.Description)

	err = be.Apply(context.Background(), rs)
	ExitOnError(err, "Apply snapshot")

	fmt.Printf("%s Successfully rolled back to %s\n", Colorize("✓", colorGreen), id)
}

package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/config"
)

type AuditCommand struct{}

func (c *AuditCommand) Name() string        { return "audit" }
func (c *AuditCommand) Description() string { return "View audit log" }

func (c *AuditCommand) Run(args []string) {
	if len(args) == 0 {
		c.usage()
		os.Exit(1)
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "list":
		c.runList(subArgs)
	case "verify":
		c.runVerify(subArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown audit subcommand: %s\n", subcommand)
		c.usage()
		os.Exit(1)
	}
}

func (c *AuditCommand) usage() {
	fmt.Fprintf(os.Stderr, "Usage: rampart audit <subcommand> [arguments]\n\n")
	fmt.Fprintf(os.Stderr, "Subcommands:\n")
	fmt.Fprintf(os.Stderr, "  %-12s %s\n", "list", "List audit events")
	fmt.Fprintf(os.Stderr, "  %-12s %s\n", "verify", "Verify audit log integrity")
}

func (c *AuditCommand) runList(args []string) {
	fs := flag.NewFlagSet("audit list", flag.ExitOnError)
	limit := fs.Int("limit", 20, "Number of events to show")
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		cfg = config.DefaultConfig()
	}

	auditStore, err := audit.NewStore(cfg.Audit)
	ExitOnError(err, "Initialize audit store")

	events, _, err := auditStore.Search(audit.AuditQuery{
		Limit: *limit,
	})
	ExitOnError(err, "Query audit log")

	fmt.Printf("%-20s %-15s %-20s %s\n", "Timestamp", "Action", "Actor", "Result")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, e := range events {
		fmt.Printf("%-20s %-15s %-20s %s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Action,
			e.Actor.Identity,
			e.Result.Status)
	}
}

func (c *AuditCommand) runVerify(args []string) {
	fs := flag.NewFlagSet("audit verify", flag.ExitOnError)
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		cfg = config.DefaultConfig()
	}

	auditStore, err := audit.NewStore(cfg.Audit)
	ExitOnError(err, "Initialize audit store")

	fmt.Printf("Verifying audit log integrity in %s...\n", cfg.Audit.Directory)
	valid, err := auditStore.VerifyIntegrity()
	if err != nil {
		ExitOnError(err, "Integrity verification failed")
	}

	if valid {
		fmt.Printf("%s Audit log integrity verified successfully.\n", Colorize("OK", colorGreen))
	} else {
		fmt.Printf("%s Audit log integrity check failed!\n", Colorize("ERROR", colorRed))
		os.Exit(1)
	}
}

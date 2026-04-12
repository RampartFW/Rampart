package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/rampartfw/rampart/internal/backend"
)

type RulesCommand struct{}

func (c *RulesCommand) Name() string        { return "rules" }
func (c *RulesCommand) Description() string { return "Manage rules directly" }

func (c *RulesCommand) Run(args []string) {
	if len(args) == 0 {
		c.usage()
		return
	}

	sub := args[0]
	subArgs := args[1:]

	switch sub {
	case "list":
		c.list(subArgs)
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

func (c *RulesCommand) usage() {
	fmt.Println("Usage: rampart rules <subcommand> [arguments]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  list      List current rules in the backend")
	fmt.Println("  delete    Delete a rule (stub)")
}

func (c *RulesCommand) list(args []string) {
	fs := flag.NewFlagSet("rules list", flag.ExitOnError)
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	be, err := backend.AutoDetect()
	ExitOnError(err, "Auto-detect backend")

	current, err := be.CurrentState()
	ExitOnError(err, "Get current state")

	if Output == "json" {
		// Mock JSON output for now
		fmt.Printf(`{"rules": %d}`, len(current.Rules))
		return
	}

	fmt.Printf("Current rules in %s:\n", be.Name())
	for _, rule := range current.Rules {
		fmt.Printf("- ID: %s, Action: %s, Protocol: %v, Source: %v, Dest: %v\n",
			rule.ID, rule.Action, rule.Match.Protocols, rule.Match.SourceNets, rule.Match.DestNets)
	}
}

func (c *RulesCommand) delete(args []string) {
	fmt.Println("Rule deletion not implemented yet.")
}

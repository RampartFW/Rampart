package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/model"
)

type RulesCommand struct{}

func (c *RulesCommand) Name() string        { return "rules" }
func (c *RulesCommand) Description() string { return "Manage rules directly (list, add, delete, stats)" }

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
	case "add":
		c.add(subArgs)
	case "delete", "remove":
		c.delete(subArgs)
	case "stats":
		c.stats(subArgs)
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
	fmt.Println("  list      List all active firewall rules")
	fmt.Println("  add       Add a new rule (temporary or permanent)")
	fmt.Println("  delete    Remove a rule by ID or name")
	fmt.Println("  stats     Show real-time packet/byte counters")
}

func (c *RulesCommand) getBackend() (backend.Backend, *config.Config) {
	cfg, _ := config.LoadConfig(ConfigPath)
	var be backend.Backend
	var err error

	if cfg != nil && cfg.Backend.Type != "auto" {
		// Manual conversion as we did in serve.go
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
			// No special settings needed for mock
		}
		be, err = backend.NewBackend(cfg.Backend.Type, bcfg)
	} else {
		be, err = backend.AutoDetect()
	}

	ExitOnError(err, "Initialize backend")
	return be, cfg
}

func (c *RulesCommand) list(args []string) {
	be, _ := c.getBackend()
	defer be.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	current, err := be.CurrentState(ctx)
	ExitOnError(err, "Fetch current rules")

	if Output == "json" {
		data, _ := json.MarshalIndent(current.Rules, "", "  ")
		fmt.Println(string(data))
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tPRIO\tPROTO\tSOURCE\tDEST\tACTION")
	
	// Sort by priority
	sort.Slice(current.Rules, func(i, j int) bool {
		return current.Rules[i].Priority < current.Rules[j].Priority
	})

	for _, r := range current.Rules {
		proto := "any"
		if len(r.Match.Protocols) > 0 {
			proto = r.Match.Protocols[0].String()
		}
		
		src := "any"
		if len(r.Match.SourceNets) > 0 {
			src = r.Match.SourceNets[0].String()
		}
		
		dst := "any"
		if len(r.Match.DestNets) > 0 {
			dst = r.Match.DestNets[0].String()
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			r.ID[:8], r.Name, r.Priority, proto, src, dst, r.Action)
	}
	w.Flush()
}

func (c *RulesCommand) add(args []string) {
	fs := flag.NewFlagSet("rules add", flag.ExitOnError)
	name := fs.String("name", "", "Rule name")
	prio := fs.Int("priority", 500, "Rule priority (0-999)")
	proto := fs.String("proto", "any", "Protocol (tcp, udp, icmp, any)")
	src := fs.String("source", "0.0.0.0/0", "Source CIDR")
	dst := fs.String("dest", "0.0.0.0/0", "Destination CIDR")
	dport := fs.Int("dport", 0, "Destination port")
	action := fs.String("action", "accept", "Action (accept, drop, reject)")
	fs.Parse(args)

	if *name == "" {
		ExitOnError(fmt.Errorf("--name is required"), "Validate flags")
	}

	be, _ := c.getBackend()
	defer be.Close()

	// 1. Get current rules
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	current, err := be.CurrentState(ctx)
	ExitOnError(err, "Get current state")

	// 2. Build new rule
	_, srcNet, _ := net.ParseCIDR(*src)
	_, dstNet, _ := net.ParseCIDR(*dst)

	newRule := model.CompiledRule{
		ID:       model.GenerateUUIDv7(),
		Name:     *name,
		Priority: *prio,
		Action:   model.ActionFromString(*action),
		Match: model.CompiledMatch{
			Protocols:  []model.Protocol{model.ProtocolFromString(*proto)},
			SourceNets: []net.IPNet{*srcNet},
			DestNets:   []net.IPNet{*dstNet},
		},
	}
	
	if *dport > 0 {
		newRule.Match.DestPorts = []model.PortRange{{Start: uint16(*dport), End: uint16(*dport)}}
	}

	// 3. Add to ruleset
	current.Rules = append(current.Rules, newRule)
	current.CompiledAt = time.Now()

	// 4. Apply
	err = be.Apply(ctx, current)
	ExitOnError(err, "Apply updated ruleset")

	fmt.Printf("Rule %q added successfully (ID: %s)\n", *name, newRule.ID)
}

func (c *RulesCommand) delete(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: rampart rules delete <ID or Name>")
		os.Exit(1)
	}
	target := args[0]

	be, _ := c.getBackend()
	defer be.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	current, err := be.CurrentState(ctx)
	ExitOnError(err, "Get current state")

	// Filter out the rule
	var newRules []model.CompiledRule
	found := false
	for _, r := range current.Rules {
		if r.ID == target || r.Name == target || strings.HasPrefix(r.ID, target) {
			found = true
			continue
		}
		newRules = append(newRules, r)
	}

	if !found {
		ExitOnError(fmt.Errorf("rule %q not found", target), "Search rule")
	}

	current.Rules = newRules
	err = be.Apply(ctx, current)
	ExitOnError(err, "Apply updated ruleset")

	fmt.Printf("Rule %q removed successfully\n", target)
}

func (c *RulesCommand) stats(args []string) {
	be, _ := c.getBackend()
	defer be.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	stats, err := be.Stats(ctx)
	ExitOnError(err, "Get stats")

	if Output == "json" {
		data, _ := json.MarshalIndent(stats, "", "  ")
		fmt.Println(string(data))
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RULE ID\tPACKETS\tBYTES")
	for id, s := range stats {
		fmt.Fprintf(w, "%s\t%d\t%d\n", id[:8], s.Packets, s.Bytes)
	}
	w.Flush()
}

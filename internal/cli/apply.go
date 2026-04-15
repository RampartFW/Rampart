package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
	"github.com/rampartfw/rampart/internal/snapshot"
)

type ApplyCommand struct{}

func (c *ApplyCommand) Name() string        { return "apply" }
func (c *ApplyCommand) Description() string { return "Apply policy from YAML file" }

func (c *ApplyCommand) getBackend(cfg *config.Config) (backend.Backend, error) {
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
		// No special settings needed
	}
	return backend.NewBackend(cfg.Backend.Type, bcfg)
}

func (c *ApplyCommand) Run(args []string) {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	file := fs.String("f", "", "Policy YAML file path")
	autoApprove := fs.Bool("auto-approve", false, "Skip confirmation prompt")
	dryRun := fs.Bool("dry-run", false, "Show plan without applying")
	fs.Parse(args)

	if *file == "" {
		fmt.Fprintln(os.Stderr, Colorize("Error", colorRed)+": -f flag is required")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		cfg = config.DefaultConfig()
	}
	vars := LoadVars()

	// 1. Parse and compile
	ps, err := engine.ParsePolicyFileWithVars(*file, vars)
	ExitOnError(err, "Parse policy")

	err = engine.ResolveIncludes(ps, *file)
	ExitOnError(err, "Resolve includes")

	compiled, err := engine.Compile(ps, vars)
	ExitOnError(err, "Compile policy")

	// 2. Conflict detection
	conflicts := engine.DetectConflicts(compiled.Rules)
	if len(conflicts) > 0 {
		fmt.Println(engine.FormatConflicts(conflicts, Output))
		hasError := false
		for _, con := range conflicts {
			if con.Severity == model.SeverityError {
				hasError = true
				break
			}
		}
		if hasError {
			os.Exit(1)
		}
	}

	// 3. Initialize backend
	be, err := c.getBackend(cfg)
	ExitOnError(err, "Initialize backend")
	defer be.Close()

	// 4. Get current state and generate plan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	current, err := be.CurrentState(ctx)
	ExitOnError(err, "Get current state")

	plan := engine.GeneratePlan(current, compiled)

	// 5. Show plan
	if Output == "json" {
		planJSON, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(planJSON))
	} else {
		c.showPlan(plan)
	}

	if *dryRun || plan.IsEmpty() {
		if plan.IsEmpty() {
			fmt.Println("\nNo changes needed. Infrastructure matches policy.")
		}
		return
	}

	// 6. Confirm
	if !*autoApprove {
		if !Confirm("Apply these changes?") {
			fmt.Println("Cancelled.")
			return
		}
	}

	// 7. Create pre-apply snapshot
	snapStore, err := snapshot.NewStore(cfg.Snapshots.Directory)
	ExitOnError(err, "Initialize snapshot store")
	
	_, err = snapStore.Create("pre-apply", "Before applying: "+*file, current)
	ExitOnError(err, "Create snapshot")

	// 8. Apply
	err = be.Apply(ctx, compiled)
	ExitOnError(err, "Apply rules")

	// 9. Audit
	auditStore, err := audit.NewStore(cfg.Audit.Directory, cfg.Audit.Retention)
	if err == nil {
		_ = auditStore.Record(model.AuditEvent{
			Action: model.AuditApply,
			Actor: model.AuditActor{
				Type:     "user",
				Identity: os.Getenv("USER"),
			},
			Result: model.AuditResult{Status: "success"},
		})
	}

	fmt.Printf("\n%s Applied %d rules (%d added, %d removed, %d modified)\n",
		Colorize("✓", colorGreen), len(compiled.Rules), plan.AddCount, plan.RemoveCount, plan.ModifyCount)
}

func (c *ApplyCommand) showPlan(plan *model.ExecutionPlan) {
	if plan.IsEmpty() {
		return
	}

	fmt.Println(Colorize("\nExecution Plan:", colorBold))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ACTION\tNAME\tPRIORITY\tDETAILS")

	for _, r := range plan.ToAdd {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			Colorize("+", colorGreen), r.Name, r.Priority, c.getRuleSummary(r))
	}

	for _, r := range plan.ToRemove {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			Colorize("-", colorRed), r.Name, r.Priority, c.getRuleSummary(r))
	}

	for _, m := range plan.ToModify {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s (changes: %s)\n",
			Colorize("~", colorYellow), m.After.Name, m.After.Priority, 
			c.getRuleSummary(m.After), strings.Join(m.Fields, ", "))
	}
	w.Flush()
}

func (c *ApplyCommand) getRuleSummary(r model.CompiledRule) string {
	proto := "any"
	if len(r.Match.Protocols) > 0 {
		proto = r.Match.Protocols[0].String()
	}
	
	dst := "any"
	if len(r.Match.DestNets) > 0 {
		dst = r.Match.DestNets[0].String()
	}

	return fmt.Sprintf("%s %s -> %s", strings.ToUpper(proto), dst, r.Action)
}

package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
	"github.com/rampartfw/rampart/internal/snapshot"
)

type ApplyCommand struct{}

func (c *ApplyCommand) Name() string        { return "apply" }
func (c *ApplyCommand) Description() string { return "Apply policy from YAML file" }

func (c *ApplyCommand) Run(args []string) {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	file := fs.String("f", "", "Policy YAML file path")
	autoApprove := fs.Bool("auto-approve", false, "Skip confirmation prompt")
	dryRun := fs.Bool("dry-run", false, "Show plan without applying")
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	if *file == "" {
		fmt.Fprintln(os.Stderr, Colorize("Error", colorRed)+": -f flag is required")
		os.Exit(1)
	}

	// 1. Parse and compile
	ps, err := engine.ParsePolicyFileWithVars(*file, LoadVars())
	ExitOnError(err, "Parse policy")

	err = engine.ResolveIncludes(ps, *file)
	ExitOnError(err, "Resolve includes")

	compiled, err := engine.Compile(ps, LoadVars())
	ExitOnError(err, "Compile policy")

	// 2. Conflict detection
	conflicts := engine.DetectConflicts(compiled.Rules)
	if len(conflicts) > 0 {
		fmt.Println(engine.FormatConflicts(conflicts, Output))
		hasError := false
		for _, c := range conflicts {
			if c.Severity == model.SeverityError {
				hasError = true
				break
			}
		}
		if hasError {
			os.Exit(1)
		}
	}

	// 3. Get current state
	be, err := backend.AutoDetect()
	ExitOnError(err, "Auto-detect backend")

	current, err := be.CurrentState()
	ExitOnError(err, "Get current state")

	// 4. Generate plan
	plan := engine.GeneratePlan(current, compiled)

	// 5. Show plan
	planJSON, _ := json.MarshalIndent(plan, "", "  ")
	fmt.Println(string(planJSON))

	if *dryRun {
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
	// For now use default directory if not provided in config
	snapDir := "snapshots"
	if _, err := os.Stat(snapDir); os.IsNotExist(err) {
		os.MkdirAll(snapDir, 0755)
	}
	snapStore, err := snapshot.NewStore(snapDir)
	ExitOnError(err, "Initialize snapshot store")
	_, err = snapStore.Create("pre-apply", "Pre-apply: "+*file, current)
	ExitOnError(err, "Create snapshot")

	// 8. Apply
	err = be.Apply(compiled)
	ExitOnError(err, "Apply rules")

	// 9. Audit
	auditDir := "audit"
	if _, err := os.Stat(auditDir); os.IsNotExist(err) {
		os.MkdirAll(auditDir, 0755)
	}
	auditStore, err := audit.NewStore(auditDir, time.Hour*2160)
	ExitOnError(err, "Initialize audit store")

	err = auditStore.Record(model.AuditEvent{
		Action: model.AuditApply,
		Actor: model.AuditActor{
			Type:     "user",
			Identity: os.Getenv("USER"),
		},
		Result: model.AuditResult{Status: "success"},
	})
	ExitOnError(err, "Record audit event")

	fmt.Printf("\n%s Applied %d rules (%d added, %d removed, %d modified)\n",
		Colorize("✓", colorGreen), len(compiled.Rules), plan.AddCount, plan.RemoveCount, plan.ModifyCount)
}

package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

type PlanCommand struct{}

func (c *PlanCommand) Name() string        { return "plan" }
func (c *PlanCommand) Description() string { return "Show execution plan without applying" }

func (c *PlanCommand) Run(args []string) {
	fs := flag.NewFlagSet("plan", flag.ExitOnError)
	file := fs.String("f", "", "Policy YAML file path")
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

	current, err := be.CurrentState(context.Background())
	ExitOnError(err, "Get current state")

	// 4. Generate plan
	plan := engine.GeneratePlan(current, compiled)

	// 5. Show plan
	planJSON, _ := json.MarshalIndent(plan, "", "  ")
	fmt.Println(string(planJSON))
}

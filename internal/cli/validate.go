package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

type ValidateCommand struct{}

func (c *ValidateCommand) Name() string        { return "validate" }
func (c *ValidateCommand) Description() string { return "Validate YAML policy file" }

func (c *ValidateCommand) Run(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
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

	fmt.Printf("%s Policy is valid.\n", Colorize("✓", colorGreen))
}

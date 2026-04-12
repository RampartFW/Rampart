package cli

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type FmtCommand struct{}

func (c *FmtCommand) Name() string        { return "fmt" }
func (c *FmtCommand) Description() string { return "Format YAML policy file" }

func (c *FmtCommand) Run(args []string) {
	fs := flag.NewFlagSet("fmt", flag.ExitOnError)
	file := fs.String("f", "", "Policy YAML file path")
	write := fs.Bool("w", false, "Write result to (source) file instead of stdout")
	RegisterGlobalFlags(fs)
	fs.Parse(args)

	if *file == "" {
		fmt.Fprintln(os.Stderr, Colorize("Error", colorRed)+": -f flag is required")
		os.Exit(1)
	}

	data, err := os.ReadFile(*file)
	ExitOnError(err, "Read policy file")

	var node yaml.Node
	err = yaml.Unmarshal(data, &node)
	ExitOnError(err, "Unmarshal policy")

	// Marshal back to string with consistent indentation
	formatted, err := yaml.Marshal(&node)
	ExitOnError(err, "Marshal policy")

	if *write {
		err = os.WriteFile(*file, formatted, 0644)
		ExitOnError(err, "Write policy file")
		fmt.Printf("%s Formatted %s\n", Colorize("✓", colorGreen), *file)
	} else {
		fmt.Print(string(formatted))
	}
}

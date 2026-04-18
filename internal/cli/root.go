package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/engine"

	_ "github.com/rampartfw/rampart/internal/backend/aws"
	_ "github.com/rampartfw/rampart/internal/backend/azure"
	_ "github.com/rampartfw/rampart/internal/backend/ebpf"
	_ "github.com/rampartfw/rampart/internal/backend/gcp"
	_ "github.com/rampartfw/rampart/internal/backend/iptables"
	_ "github.com/rampartfw/rampart/internal/backend/nftables"
)

// Global flags
var (
	ConfigPath string = "rampart.yaml"
	Output     string
	Verbose    bool
	NoColor    bool
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

// RegisterGlobalFlags adds the global flags to the provided FlagSet.
func RegisterGlobalFlags(fs *flag.FlagSet) {
	fs.StringVar(&ConfigPath, "config", "", "Path to config file")
	fs.StringVar(&Output, "output", "text", "Output format (text, json)")
	fs.StringVar(&Output, "o", "text", "Output format (text, json) (shorthand)")
	fs.BoolVar(&Verbose, "verbose", false, "Enable verbose output")
	fs.BoolVar(&NoColor, "no-color", false, "Disable color output")
}

// Colorize returns the text wrapped in the specified color if NoColor is false.
func Colorize(text, color string) string {
	if NoColor || os.Getenv("NO_COLOR") != "" {
		return text
	}
	return color + text + colorReset
}

// ExitOnError prints the error and context and exits the program.
func ExitOnError(err error, context string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", Colorize("Error", colorRed), fmt.Sprintf("%s: %v", context, err))
		os.Exit(1)
	}
}

// Confirm asks for user confirmation and returns true if accepted.
func Confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// Subcommand is an interface for all CLI subcommands.
type Subcommand interface {
	Name() string
	Description() string
	Run(args []string)
}

// LoadConfig loads the Rampart configuration.
func LoadConfig() *config.Config {
	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to load config: %v\n", Colorize("Warning", colorYellow), err)
		return config.DefaultConfig()
	}
	return cfg
}

// LoadVars loads the policy variables.
func LoadVars() map[string]interface{} {
	// 1. Try default location
	vars, err := engine.ParseVariablesFile("rampart-vars.yaml")
	if err == nil {
		return vars
	}

	// 2. If not found, look in ~/.config/rampart/vars.yaml
	home, _ := os.UserHomeDir()
	vars, err = engine.ParseVariablesFile(home + "/.config/rampart/vars.yaml")
	if err == nil {
		return vars
	}

	return make(map[string]interface{})
}

var subcommands = []Subcommand{
	&ApplyCommand{},
	&AuditCommand{},
	&CertCommand{},
	&ClusterCommand{},
	&FmtCommand{},
	&PlanCommand{},
	&RollbackCommand{},
	&RulesCommand{},
	&ServeCommand{},
	&SnapshotCommand{},
	&ValidateCommand{},
	&VersionCommand{},
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: rampart [global flags] <command> [arguments]\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	for _, cmd := range subcommands {
		fmt.Fprintf(os.Stderr, "  %-12s %s\n", cmd.Name(), cmd.Description())
	}
	fmt.Fprintf(os.Stderr, "\nGlobal Flags:\n")
	fs := flag.NewFlagSet("global", flag.ContinueOnError)
	RegisterGlobalFlags(fs)
	fs.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nUse \"rampart <command> --help\" for more information about a command.\n")
}

func Execute() {
	// 1. Setup global flags on the default flagset
	RegisterGlobalFlags(flag.CommandLine)
	flag.Usage = usage
	flag.Parse()

	// 2. Check for command
	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	commandName := args[0]
	commandArgs := args[1:]

	if commandName == "help" {
		usage()
		os.Exit(0)
	}

	for _, cmd := range subcommands {
		if cmd.Name() == commandName {
			cmd.Run(commandArgs)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", commandName)
	usage()
	os.Exit(1)
}

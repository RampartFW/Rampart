package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/rampartfw/rampart/internal/cert"
)

type CertCommand struct{}

func (c *CertCommand) Name() string        { return "cert" }
func (c *CertCommand) Description() string { return "Manage cluster certificates" }

func (c *CertCommand) Run(args []string) {
	if len(args) == 0 {
		c.usage()
		return
	}

	sub := args[0]
	subArgs := args[1:]

	switch sub {
	case "init":
		c.runInit(subArgs)
	case "generate":
		c.runGenerate(subArgs)
	case "help", "-h", "--help":
		c.usage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", sub)
		c.usage()
		os.Exit(1)
	}
}

func (c *CertCommand) runInit(args []string) {
	fs := flag.NewFlagSet("cert init", flag.ExitOnError)
	caDir := fs.String("ca-dir", "certs", "Directory to store CA certificates")
	fs.Parse(args)

	if err := cert.InitCA(*caDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully initialized CA in %s\n", *caDir)
}

func (c *CertCommand) runGenerate(args []string) {
	fs := flag.NewFlagSet("cert generate", flag.ExitOnError)
	nodeName := fs.String("node-name", "", "Name of the node")
	caDir := fs.String("ca-dir", "certs", "Directory containing CA certificates")
	fs.Parse(args)

	if *nodeName == "" {
		fmt.Fprintln(os.Stderr, "Error: --node-name is required")
		fs.Usage()
		os.Exit(1)
	}

	if err := cert.GenerateNodeCert(*nodeName, *caDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated certificate for node %s in %s\n", *nodeName, *caDir)
}

func (c *CertCommand) usage() {
	fmt.Println("Usage: rampart cert <subcommand> [arguments]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  init        Initialize cluster CA")
	fmt.Println("  generate    Generate node certificate")
	fmt.Println("\nArguments for init:")
	fmt.Println("  --ca-dir    Directory to store CA certificates (default: certs)")
	fmt.Println("\nArguments for generate:")
	fmt.Println("  --node-name Name of the node (required)")
	fmt.Println("  --ca-dir    Directory containing CA certificates (default: certs)")
}

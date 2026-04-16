package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/model"
)

type ClusterCommand struct{}

func (c *ClusterCommand) Name() string        { return "cluster" }
func (c *ClusterCommand) Description() string { return "Manage and monitor cluster status" }

func (c *ClusterCommand) Run(args []string) {
	if len(args) == 0 {
		c.usage()
		return
	}

	sub := args[0]
	subArgs := args[1:]

	switch sub {
	case "status":
		c.runStatus(subArgs)
	case "init", "join", "leave", "elect":
		fmt.Printf("Subcommand %s is currently managed via configuration and 'rampart serve'.\n", sub)
	case "help", "-h", "--help":
		c.usage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", sub)
		c.usage()
		os.Exit(1)
	}
}

func (c *ClusterCommand) runStatus(args []string) {
	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		cfg = config.DefaultConfig()
	}

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://%s/api/v1/cluster/status", cfg.Server.Listen)
	
	req, _ := http.NewRequest("GET", url, nil)
	// Add default test key if present in config for auth
	if len(cfg.API.Keys) > 0 {
		req.Header.Set("Authorization", "Bearer "+cfg.API.Keys[0].Key)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to Rampart API at %s. Is the server running?\n", url)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct{ Message string } `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Fprintf(os.Stderr, "Error: API returned %d - %s\n", resp.StatusCode, errResp.Error.Message)
		os.Exit(1)
	}

	var apiResp struct {
		Status string           `json:"status"`
		Data   model.NodeStatus `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse API response: %v\n", err)
		os.Exit(1)
	}

	status := apiResp.Data
	fmt.Println(Colorize("\nCluster Node Status:", colorBold))
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NODE ID\tSTATE\tHEALTH\tLAST SYNC")
	
	healthStr := Colorize("✓ Healthy", colorGreen)
	if !status.IsHealthy {
		healthStr = Colorize("✗ Unhealthy", colorRed)
	}

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		status.ID,
		strings.ToUpper(string(status.State)),
		healthStr,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	w.Flush()
	fmt.Println()
}

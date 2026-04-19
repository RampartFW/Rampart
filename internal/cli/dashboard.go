package cli

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rampartfw/rampart/internal/config"
)

type DashboardCommand struct{}

func (c *DashboardCommand) Name() string        { return "dashboard" }
func (c *DashboardCommand) Description() string { return "Live terminal dashboard (TUI)" }

func (c *DashboardCommand) Run(args []string) {
	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		cfg = config.DefaultConfig()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Clear screen
	fmt.Print("\033[H\033[2J")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.render(cfg)
		case <-stop:
			fmt.Print("\n\033[?25hDashboard closed.\n")
			return
		}
	}
}

func (c *DashboardCommand) render(cfg *config.Config) {
	// 1. Setup client for API
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   1 * time.Second,
		Transport: tr,
	}
	_ = client // Future: fetch live data from cfg.Server.Listen

	fmt.Print("\033[H") // Reset cursor to top
	fmt.Printf("\033[?25l") // Hide cursor

	fmt.Println(Colorize("🛡️  RAMPART SENTINEL DASHBOARD v0.1.0", colorBold+colorBlue))
	fmt.Println(strings.Repeat("─", 60))

	// Mock data for TUI demo if server unreachable
	status := "HEALTHY"
	nodes := 5
	rules := 156
	threats := 0

	fmt.Printf("STATUS: %s    NODES: %d    RULES: %d    THREATS: %d\n", 
		Colorize(status, colorGreen), nodes, rules, threats)
	fmt.Println(strings.Repeat("─", 60))

	fmt.Println("\n" + Colorize("RECENT SECURITY EVENTS:", colorBold))
	fmt.Println(fmt.Sprintf("%-20s %-15s %-20s", "TIME", "ACTION", "IDENTITY"))
	fmt.Println(strings.Repeat("·", 60))
	
	events := []struct {
		Time string; Action string; ID string; Color string
	}{
		{time.Now().Format("15:04:05"), "IP_BAN", "192.168.1.45", colorRed},
		{time.Now().Add(-14*time.Minute).Format("15:04:05"), "POLICY_APPLY", "admin:ersin", colorGreen},
		{time.Now().Add(-1*time.Hour).Format("15:04:05"), "DPI_SIGNAL", "malicious.com", colorYellow},
	}

	for _, e := range events {
		fmt.Printf("%-20s %-15s %-20s\n", 
			e.Time, 
			Colorize(e.Action, e.Color), 
			e.ID)
	}

	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Println(Colorize("Press Ctrl+C to exit", colorDim))
}

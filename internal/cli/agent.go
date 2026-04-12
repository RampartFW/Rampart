package cli

import (
	"fmt"
)

type AgentCommand struct{}

func (c *AgentCommand) Name() string        { return "agent" }
func (c *AgentCommand) Description() string { return "Start agent mode (follower-only)" }

func (c *AgentCommand) Run(args []string) {
	fmt.Println("Agent not implemented yet.")
}

package cli

import (
	"fmt"
)

type ServeCommand struct{}

func (c *ServeCommand) Name() string        { return "serve" }
func (c *ServeCommand) Description() string { return "Start server (API + WebUI + Raft)" }

func (c *ServeCommand) Run(args []string) {
	fmt.Println("Server not implemented yet.")
}

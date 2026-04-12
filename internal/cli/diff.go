package cli

import (
	"fmt"
)

type DiffCommand struct{}

func (c *DiffCommand) Name() string        { return "diff" }
func (c *DiffCommand) Description() string { return "Diff two policy files" }

func (c *DiffCommand) Run(args []string) {
	fmt.Println("Policy diff not implemented yet.")
}

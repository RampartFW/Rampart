package cli

import (
	"fmt"
)

type ImportCommand struct{}

func (c *ImportCommand) Name() string        { return "import" }
func (c *ImportCommand) Description() string { return "Import existing rules" }

func (c *ImportCommand) Run(args []string) {
	fmt.Println("Rule import not implemented yet.")
}

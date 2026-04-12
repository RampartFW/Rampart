package cli

import (
	"fmt"
)

type ExportCommand struct{}

func (c *ExportCommand) Name() string        { return "export" }
func (c *ExportCommand) Description() string { return "Export current rules" }

func (c *ExportCommand) Run(args []string) {
	fmt.Println("Rule export not implemented yet.")
}

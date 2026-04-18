package cli

import (
	"fmt"
)

type VersionCommand struct{}

func (c *VersionCommand) Name() string        { return "version" }
func (c *VersionCommand) Description() string { return "Show version info" }

func (c *VersionCommand) Run(args []string) {
	fmt.Println("rampart version 0.1.0")
}

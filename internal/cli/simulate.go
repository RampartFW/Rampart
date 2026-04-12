package cli

import (
	"fmt"
)

type SimulateCommand struct{}

func (c *SimulateCommand) Name() string        { return "simulate" }
func (c *SimulateCommand) Description() string { return "Simulate a packet" }

func (c *SimulateCommand) Run(args []string) {
	fmt.Println("Packet simulation not implemented yet.")
}

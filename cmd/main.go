package main

import (
	"context"
	"log"

	"github.com/bogatyr285/sumcalc-testtask/cmd/commands"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()

	c := cobra.Command{}
	c.AddCommand(commands.NewServeCmd())
	// other commands here

	if err := c.ExecuteContext(ctx); err != nil {
		log.Fatalf("cmd err: %v", err)
	}
}

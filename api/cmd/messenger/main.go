package main

import (
	"os"

	"github.com/and-period/furumaru/api/internal/messenger/cmd"
	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{Use: "messenger [command]"}
	cmd.RegisterCommand(c)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

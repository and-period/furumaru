package main

import (
	"os"

	"github.com/and-period/furumaru/api/internal/media/cmd"
	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{Use: "media [command]"}
	cmd.RegisterCommand(c)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

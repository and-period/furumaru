package cmd

import (
	"github.com/and-period/furumaru/api/internal/messenger/cmd/scheduler"
	"github.com/and-period/furumaru/api/internal/messenger/cmd/worker"
	"github.com/spf13/cobra"
)

func RegisterCommand(registry *cobra.Command) {
	registry.AddCommand(
		scheduler.NewApp().Command,
		worker.NewApp().Command,
	)
}

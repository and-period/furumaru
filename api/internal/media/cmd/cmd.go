package cmd

import (
	"github.com/and-period/furumaru/api/internal/media/cmd/resizer"
	"github.com/and-period/furumaru/api/internal/media/cmd/scheduler"
	"github.com/and-period/furumaru/api/internal/media/cmd/updater"
	"github.com/and-period/furumaru/api/internal/media/cmd/uploader"
	"github.com/spf13/cobra"
)

func RegisterCommand(registry *cobra.Command) {
	registry.AddCommand(
		resizer.NewApp().Command,
		scheduler.NewApp().Command,
		updater.NewApp().Command,
		uploader.NewApp().Command,
	)
}

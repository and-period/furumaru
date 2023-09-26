package cmd

import (
	"github.com/and-period/furumaru/api/internal/gateway/cmd/admin"
	"github.com/and-period/furumaru/api/internal/gateway/cmd/user"
	"github.com/spf13/cobra"
)

func RegisterCommand(registry *cobra.Command) {
	registry.AddCommand(
		admin.NewApp().Command,
		user.NewApp().Command,
	)
}

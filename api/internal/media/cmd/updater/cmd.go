package updater

import "github.com/spf13/cobra"

type app struct {
	*cobra.Command
}

//nolint:revive
func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "updater",
		Short: "media updater",
	}
	app := &app{Command: cmd}
	app.RunE = func(c *cobra.Command, args []string) error {
		return app.run()
	}
	return app
}

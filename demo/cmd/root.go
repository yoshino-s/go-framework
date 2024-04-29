package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/telemetry"
)

var App = application.NewMainApplication()

var (
	telemetry_app = telemetry.New()
	rootCmd       = &cobra.Command{
		Use: "demo",
	}
)

func init() {
	cobra.OnInitialize(func() {
		configuration.Setup("demo")

		App.Append(telemetry_app)
	})

	App.Configuration().Register(rootCmd.PersistentFlags())
	telemetry_app.Configuration().Register(rootCmd.PersistentFlags())
	configuration.GenerateConfiguration.Register(rootCmd.PersistentFlags())

	rootCmd.AddCommand(common.VersionCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

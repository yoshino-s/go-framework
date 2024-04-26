package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

var App = application.NewMainApplication()

var rootCmd = &cobra.Command{
	Use: "demo",
}

func init() {
	cobra.OnInitialize(func() {
		configuration.Setup("demo")
	})

	configuration.LogConfiguration.Register(rootCmd.PersistentFlags())
	configuration.GenerateConfiguration.Register(rootCmd.PersistentFlags())
	configuration.TelemetryConfiguration.Register(rootCmd.PersistentFlags())

	rootCmd.AddCommand(common.VersionCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

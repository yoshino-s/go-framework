package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/cmd"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

var (
	App     = application.NewMainApplication()
	rootCmd *cobra.Command
)

func init() {
	common.AppName = "demo"

	rootCmd = &cobra.Command{
		Use: common.AppName,
	}

	cobra.OnInitialize(func() {
		configuration.Setup(common.AppName)
	})

	App.Configuration().Register(rootCmd.PersistentFlags())
	configuration.GenerateConfiguration.Register(rootCmd.PersistentFlags())

	rootCmd.AddCommand(cmd.VersionCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

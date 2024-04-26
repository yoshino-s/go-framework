package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/handlers/http"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		handler, err := http.New(
			configuration.HttpHandlerConfiguration.Config,
		)
		if err != nil {
			panic(err)
		}
		App.Add(handler)
		App.Go(context.TODO())
	},
}

func init() {
	configuration.HttpHandlerConfiguration.Register(serveCmd.Flags())
	rootCmd.AddCommand(serveCmd)
}

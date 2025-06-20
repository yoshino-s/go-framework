package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yoshino-s/go-framework/internal/demo/app"
)

var (
	demoApp  = app.New()
	serveCmd = &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			App.Append(demoApp)
			App.Append(app.NewService())
			App.Go(context.TODO())
		},
	}
)

func init() {
	demoApp.Configuration().Register(serveCmd.Flags())
	rootCmd.AddCommand(serveCmd)
}

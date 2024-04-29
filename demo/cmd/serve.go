package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yoshino-s/go-framework/demo/app"
	"github.com/yoshino-s/go-framework/handlers/http"
)

var (
	handler  = http.New()
	demoApp  = app.New()
	serveCmd = &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			App.SubApplication.Append(handler)
			App.SubApplication.Append(demoApp)
			App.Go(context.TODO())
		},
	}
)

func init() {
	handler.Configuration().Register(serveCmd.Flags())
	rootCmd.AddCommand(serveCmd)
}

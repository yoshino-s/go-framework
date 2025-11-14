package cmd

import (
	"context"

	"github.com/casbin/casbin/v2/model"
	"github.com/spf13/cobra"
	authz_casbin "github.com/yoshino-s/go-framework/authz/casbin"
	"github.com/yoshino-s/go-framework/internal/demo/app"
)

var (
	demoApp             = app.New()
	casbinAuthorization = app.NewSelfCasbinAuthorization()
	serveCmd            = &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			casbinModel, _ := model.NewModelFromString(authz_casbin.DefaultModelConf)
			App.Register(casbinModel)
			App.Append(casbinAuthorization)
			App.Append(demoApp)
			App.Append(app.NewService())

			App.Go(context.TODO())
		},
	}
)

func init() {
	demoApp.Configuration().Register(serveCmd.Flags())
	casbinAuthorization.Configuration().Register(serveCmd.Flags())
	rootCmd.AddCommand(serveCmd)
}

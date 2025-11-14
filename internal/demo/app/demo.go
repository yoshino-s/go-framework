package app

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/handlers/http"
)

var _ application.Application = (*DemoApp)(nil)

type DemoApp struct {
	*http.Handler
	Service             *Service                 `inject:""`
	CasbinAuthorization *SelfCasbinAuthorization `inject:""`
}

func New() *DemoApp {
	return &DemoApp{
		Handler: http.New(),
	}
}

func (a *DemoApp) Setup(ctx context.Context) {
	a.Handler.Setup(ctx)
	a.Handler.GET("/random", func(c echo.Context) error {
		fmt.Println("Received request for random number")
		if a.Service == nil {
			return c.JSON(500, "Service not initialized")
		}
		randomNumber := a.Service.GetRandomNumber()
		return c.JSON(200, map[string]int{"random_number": randomNumber})
	})

	a.Handler.Group("/perm").Any("/*", func(c echo.Context) error {
		user := c.QueryParam("user")
		obj := "/" + c.Param("*")
		act := c.Request().Method
		allowed, err := a.CasbinAuthorization.Enforce(user, obj, act)
		if err != nil {
			return c.JSON(500, "Error during authorization: "+err.Error())
		}
		if !allowed {
			return c.JSON(403, "Forbidden")
		}
		return c.JSON(200, "Access granted to "+user+" for "+obj+" with action "+act)
	})

	a.Handler.GET("/casbin-dump", func(c echo.Context) error {
		dump, err := a.CasbinAuthorization.Dump()
		if err != nil {
			return c.JSON(500, "Error retrieving dump: "+err.Error())
		}
		return c.JSON(200, dump)
	})
}

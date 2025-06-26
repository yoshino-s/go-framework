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
	Service *Service `inject:""`
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
}

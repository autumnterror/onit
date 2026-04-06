package net

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/autumnterror/onit/internal/service"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Echo struct {
	echo *echo.Echo
	repo service.Repo
}

func New(
	repo service.Repo,
) *Echo {
	e := &Echo{
		echo: echo.New(),
		repo: repo,
	}

	e.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1"), nil
		},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderContentType},
		AllowCredentials: true,
	}))

	api := e.echo.Group("/api")
	{
		api.GET("/health", e.Health)
		prod := api.Group("/product")
		{
			prod.GET("", e.GetById)
			prod.GET("/all", e.GetAll)
			prod.DELETE("", e.Delete)
			prod.PUT("", e.Update)
			prod.POST("", e.Create)
		}
	}

	return e
}

func (e *Echo) MustRun() {
	const op = "net.Run"

	if err := e.echo.Start(fmt.Sprintf(":%d", 8080)); err != nil && !errors.Is(http.ErrServerClosed, err) {
		e.echo.Logger.Fatal(format.Error(op, err))
	}
}

func (e *Echo) Stop() error {
	const op = "net.Stop"

	if err := e.echo.Close(); err != nil {
		return format.Error(op, err)
	}
	return nil
}

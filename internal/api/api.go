package api

import (
	"context"
	"fmt"

	"github.com/Marsredskies/todo-list/internal/db"
	"github.com/Marsredskies/todo-list/internal/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type API struct {
	echo *echo.Echo
	db   *db.DB
}

func MustInitNewAPI(ctx context.Context, cnf envconfig.Config) API {
	api, err := New(ctx, cnf)
	if err != nil {
		panic(fmt.Errorf("failed to initialize API: %v", err))
	}
	return api

}

func New(ctx context.Context, cnf envconfig.Config) (API, error) {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"method":"${method}","uri":"${uri}","status":${status}` + "\n",
	}))
	e.Pre(middleware.RemoveTrailingSlash())

	db, err := db.New(ctx, cnf)
	if err != nil {
		return API{}, err
	}

	api := API{
		echo: e,
		db:   db,
	}

	e.GET("/list", nil)
	e.POST("/create", nil)
	e.PATCH("/update", nil)

	return api, nil
}

func (a *API) StartServer(port int) error {
	err := a.echo.Start(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.echo.Shutdown(ctx)
}

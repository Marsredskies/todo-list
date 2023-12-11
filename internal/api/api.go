package api

import (
	"context"
	"fmt"

	"github.com/Marsredskies/todo-list/internal/db"
	"github.com/Marsredskies/todo-list/internal/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

type API struct {
	echo   *echo.Echo
	token  string
	r      Repository
	logger *log.Logger
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

	if cnf.StaticToken != "" {
		e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "query:token",
			Validator: NewTokenValidator(cnf),
		}))
	}

	dbClient, err := db.New(ctx, cnf)
	if err != nil {
		return API{}, err
	}

	logger := log.New()

	r := db.NewTaskRepo(dbClient)

	api := API{
		echo:   e,
		r:      r,
		logger: logger,
	}

	e.GET("/search-with-filters", api.handleFindTask)
	e.POST("/create", api.handleCreateTask)
	e.PATCH("/update-by-id", api.handleUpdateTask)
	e.DELETE("/delete", api.handleDeleteTask)

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

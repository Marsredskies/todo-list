package api

import (
	"context"
	"fmt"

	_ "github.com/Marsredskies/todo-list/docs"
	"github.com/Marsredskies/todo-list/internal/db"
	"github.com/Marsredskies/todo-list/internal/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type API struct {
	echo   *echo.Echo
	token  string
	r      Repository
	Logger *log.Logger
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
	e.Use(middleware.CORS())

	if cnf.StaticToken != "" {
		e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "header:token",
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
		Logger: logger,
	}

	e.GET("/search-with-filters", api.handleFindTask)
	e.POST("/create", api.handleCreateTask)
	e.PATCH("/update-by-id", api.handleUpdateTask)
	e.DELETE("/delete", api.handleDeleteTask)

	return api, nil
}

func (a *API) StartServer() error {
	err := a.echo.Start(":8080")
	if err != nil {
		return err
	}
	return nil
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.echo.Shutdown(ctx)
}

func (a *API) StartSwagger() {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8081"))
}

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Marsredskies/todo-list/internal/api"
	"github.com/Marsredskies/todo-list/internal/db"
	"github.com/Marsredskies/todo-list/internal/envconfig"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cnf := envconfig.MustGetConfig()

	db.MustApplyMigrations(ctx, cnf)

	api := api.MustInitNewAPI(ctx, cnf)

	go func() {
		if err := api.StartServer(); err != nil && err != http.ErrServerClosed {
			api.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := api.Shutdown(ctx); err != nil {
		api.Logger.Fatal(err)
	}
}

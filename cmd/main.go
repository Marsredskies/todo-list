package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Marsredskies/todo-list/internal/api"
	"github.com/Marsredskies/todo-list/internal/db"
	"github.com/Marsredskies/todo-list/internal/envconfig"
)

var GitCommitSHA string

func main() {
	log.Printf("Starting with CommitSHA: %s", GitCommitSHA)

	ctx, cancel := context.WithCancel(context.Background())

	cnf := envconfig.MustGetConfig()

	db.MustApplyMigrations(ctx, cnf)

	api := api.MustInitNewAPI(ctx, cnf)

	api.StartServer(cnf.Port)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	defer func() {
		if err := api.Shutdown(ctx); err != nil {
			log.Println("error during shutting down the main server: ", err)
		}
	}()

	<-exit
	cancel()
	log.Println("shutting down")
}

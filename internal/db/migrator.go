package db

import (
	"context"
	"fmt"

	"github.com/Marsredskies/todo-list/internal/envconfig"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var index map[int]struct{}
var allMigrations []*migrate.Migration

func init() {
	index = make(map[int]struct{})
	allMigrations = []*migrate.Migration{}
}

func ApplyMigrations(ctx context.Context, cnf envconfig.Config) error {
	db, err := ConnectDB(ctx, cnf)
	if err != nil {
		return err
	}

	m := migrate.MemoryMigrationSource{
		Migrations: allMigrations,
	}

	n, err := migrate.Exec(db.DB, "postgres", &m, migrate.Up)
	if err != nil {
		fmt.Printf("error at migration: %s", err)

		if pqErr, ok := err.(*pq.Error); ok {
			panic(fmt.Errorf("migration err: %#v", pqErr))
		}
		if txErr, ok := err.(*migrate.TxError); ok {
			panic(fmt.Errorf("migratio err: %#v", txErr.Err))
		}

		return err
	}

	if n > 0 {
		fmt.Printf("Applied %d migrations!\n", n)
	}

	err = db.Close()
	if err != nil {
		fmt.Printf("close db err %s", err)
	}

	return nil
}

func AddMigration(id int, up string) {
	if _, found := index[id]; found {
		panic(fmt.Sprintf("migration %d duplication", id))
	}
	index[id] = struct{}{}

	m := migrate.Migration{
		Id:   fmt.Sprintf("%d", id),
		Up:   []string{up},
		Down: []string{""},
	}

	allMigrations = append(allMigrations, &m)
}

func MustApplyMigrations(ctx context.Context, cnf envconfig.Config) {
	err := ApplyMigrations(ctx, cnf)
	if err != nil {
		pgErr := ExtractPqError(err)
		if pgErr != nil {
			panic(pgErr.Message)
		} else {
			panic(fmt.Errorf("migration apply: %s", err))
		}
	}
}

func ExtractPqError(err error) *pq.Error {
	txError, ok := err.(*migrate.TxError)
	if !ok {
		return nil
	}
	e, ok := txError.Err.(*pq.Error)
	if !ok {
		return nil
	}
	return e
}

func DropMigrations(db *sqlx.DB) {
	_, err := db.Exec(
		`DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public`)
	if err != nil {
		panic(fmt.Errorf("couldn't drop migrations: %v", err))
	}
}

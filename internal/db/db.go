package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Marsredskies/todo-list/envconfig"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func New(ctx context.Context, cnf envconfig.Config) (*DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db, err := ConnectDB(ctx, cnf)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func ConnectDB(ctx context.Context, cnf envconfig.Config) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	db, err := sqlx.Open("postgres", cnf.PgURL)
	if err != nil {
		log.Println("qlx.Open", err.Error())
		return nil, err
	}

	db.SetConnMaxIdleTime(10 * time.Second)
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)

	return db, nil
}

func RequireNewDBClient(ctx context.Context, cnf envconfig.Config) *DB {
	db, err := New(ctx, cnf)
	if err != nil {
		panic(fmt.Errorf("failed to initialise db client: %v", err))
	}
	return db
}

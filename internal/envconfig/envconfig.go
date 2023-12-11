package envconfig

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        int    `envconfig:"APP_PORT" default:"8080"`
	PgURL       string `envconfig:"APP_PG_URL" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	StaticToken string `envconfig:"APP_STATIC_TOKEN" default:"test_token"`
}

func GetConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("APP", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func MustGetConfig() Config {
	cnf, err := GetConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse env: %v", err))
	}
	return cnf
}

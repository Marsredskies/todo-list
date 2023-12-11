package envconfig

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port  int    `envconfig:"PORT" default:"8080"`
	PgURL string `envconfig:"PG_URL" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	Token string `envconfig:"TOKEN" default:"test_token"`
}

func GetConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("core", &cfg); err != nil {
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

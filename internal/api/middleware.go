package api

import (
	"github.com/Marsredskies/todo-list/internal/envconfig"
	"github.com/labstack/echo"
)

func NewTokenValidator(cnf envconfig.Config) func(key string, c echo.Context) (bool, error) {
	return func(key string, c echo.Context) (bool, error) {
		return key == cnf.StaticToken, nil
	}
}

package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func (a *API) handleCreateTask(c echo.Context) error {
	var params Task
	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, incorrectDataFormat)
	}

	err = params.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := a.saveTaskToDb(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("task has been created successfuly with ID %d", id))

}

func (a *API) handleUpdateTask(c echo.Context) error {
	var params Task
	err := c.Bind(&params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, incorrectDataFormat)
	}

	if params.ID == 0 {
		return c.JSON(http.StatusBadRequest, noIdProvided)
	}

	exists, err := a.checkIfTaskExists(params.ID)
	if err != nil {
		log.Println("checkIfTaskExists: ", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if !exists {
		return c.JSON(http.StatusNotFound, noTasksFound)
	}

	err = params.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = a.updateTaskById(c.Request().Context(), params)
	if err != nil {
		log.Println("updateTaskById: ", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("task has been updated successfuly"))

}

func (a *API) saveTaskToDb(ctx context.Context, params Task) (int64, error) {
	query, args, err := params.sqlInsert()
	if err != nil {
		return 0, err
	}

	return a.db.ExecCtxReturningId(ctx, query, args...)
}

func (a *API) updateTaskById(ctx context.Context, params Task) error {
	query, args, err := params.sqlUpdate()
	if err != nil {
		return err
	}

	err = a.db.ExecCtx(ctx, query, args...)
	return err
}

func (a *API) checkIfTaskExists(id int64) (bool, error) {
	var task Task
	err := a.db.GetById(&task,
		`SELECT name, description, assignee, status 
			FROM public.tasks
				WHERE id = $1 AND deleted_at IS NULL`, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

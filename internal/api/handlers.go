package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Marsredskies/todo-list/internal/models"
	"github.com/labstack/echo"
)

func (a *API) handleCreateTask(c echo.Context) error {
	var params models.Task
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
	var params models.Task
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

func (a *API) handleDeleteTask(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, invalidIdFormat)
	}

	err = a.deleteById(c.Request().Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("task with id %d has been deleted", id))
}

func (a *API) handleFindTask(c echo.Context) error {
	params := models.Task{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Assignee:    c.FormValue("assignee"),
		Status:      c.FormValue("status"),
	}

	results, err := a.getMatchingTasks(c.Request().Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

func (a *API) getMatchingTasks(ctx context.Context, params models.Task) ([]models.Task, error) {
	var results []models.Task
	query, args, err := params.SqlSelectLike()
	if err != nil {
		return nil, err
	}

	err = a.db.SelectCtx(ctx, &results, query, args...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (a *API) saveTaskToDb(ctx context.Context, params models.Task) (int64, error) {
	query, args, err := params.SqlInsert()
	if err != nil {
		return 0, err
	}

	return a.db.ExecCtxReturningId(ctx, query, args...)
}

func (a *API) updateTaskById(ctx context.Context, params models.Task) error {
	query, args, err := params.SqlUpdate()
	if err != nil {
		return err
	}

	err = a.db.ExecCtx(ctx, query, args...)
	return err
}

func (a *API) checkIfTaskExists(id int64) (bool, error) {
	var task models.Task
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

func (a *API) deleteById(ctx context.Context, id int64) error {
	err := a.db.ExecCtx(ctx,
		`UPDATE public.tasks 
			SET deleted_at = now() 
				WHERE id = $1 AND deleted_at IS NULL`, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no tasks found")
		}
		return err
	}

	return nil

}

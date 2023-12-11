package api

import (
	"fmt"
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

	id, err := a.r.SaveTask(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	a.logger.Infof("task with id %d created: %+v", id, params)

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

	exists, err := a.r.CheckIfTaskExists(params.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if !exists {
		return c.JSON(http.StatusNotFound, noTasksFound)
	}

	err = params.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = a.r.UpdateTaskById(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	a.logger.Infof("task with updated: %+v", params)

	return c.JSON(http.StatusOK, fmt.Sprintf("task has been updated successfuly"))
}

func (a *API) handleDeleteTask(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, invalidIdFormat)
	}

	err = a.r.DeleteById(c.Request().Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	a.logger.Infof("task with id %d has been deleted", id)

	return c.JSON(http.StatusOK, fmt.Sprintf("task with id %d has been deleted", id))
}

func (a *API) handleFindTask(c echo.Context) error {
	params := models.Task{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Assignee:    c.FormValue("assignee"),
		Status:      c.FormValue("status"),
	}

	results, err := a.r.GetMatchingTasks(c.Request().Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Marsredskies/todo-list/internal/models"
	"github.com/labstack/echo/v4"
)

// CreateTask creates a task with name, description, assignee and status.
//
//	@Security 	  StaticTokenAuth
//	@Summary      CreateTask
//	@Description  Creates to-do entry in the database. Name, description and assignee are required fields. If status is empty default "to do" will be set
//	@Tags         to-do list
//	@Accept       json
//	@Produce      json
//	@Param        input body models.Request true "Task description"
//	@Success      200  {object}  models.Task
//	@Failure      400,500  {string} stirng
//	@Router       /create [post]
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

	params.ID = id
	a.Logger.Infof("task with id %d created: %+v", id, params)

	return c.JSON(http.StatusOK, params)
}

// UpdateTask creates a task with name, description, assignee and status.
//
//	@Security 	  StaticTokenAuth
//	@Summary      UpdateTask
//	@Description  Update existing task by it's id
//	@Tags         to-do list
//	@Accept       json
//	@Produce      json
//	@Param        input body models.Task true "All fields are optional except the ID. Empty fileds will be ommited"
//	@Success      200  {object}  string
//	@Failure      400,500  {string}  stirng
//	@Router       /update-by-id [patch]
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

	err = params.ValidateStatus()
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	err = a.r.UpdateTaskById(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	a.Logger.Infof("task with updated: %+v", params)

	return c.JSON(http.StatusOK, fmt.Sprintf("task has been updated successfuly"))
}

// DeleteTask creates a task with name, description, assignee and status.
//
//	@Security 	  StaticTokenAuth
//	@Summary      DeleteTask
//	@Description  Delete task by it's id
//	@Tags         to-do list
//	@Accept       json
//	@Produce      json
//	@Param        id query int true "Task id"
//	@Success      200  {object}  string
//	@Failure      400,500  {string}  stirng
//	@Router       /delete [delete]
func (a *API) handleDeleteTask(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, invalidIdFormat)
	}

	exists, err := a.r.CheckIfTaskExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if !exists {
		return c.JSON(http.StatusNotFound, noTasksFound)
	}

	err = a.r.DeleteById(c.Request().Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	a.Logger.Infof("task with id %d has been deleted", id)

	return c.JSON(http.StatusOK, fmt.Sprintf("task with id %d has been deleted", id))
}

// FindTask finds all matching tasks by any mathching paramater(s)
//
//		@Security 	  StaticTokenAuth
//		@Summary      FindTask
//		@Description  Find a task by one or multiple parameters. Values may not be exact same as stored in db. It can be one search parameter or multuple 'alike' ones.
//		@Tags         to-do list
//		@Accept       json
//		@Produce      json
//		@Param		  name query string false "Name"
//		@Param		  description query string false "Description"
//		@Param		  assignee query string false "Assignee"
//	 	@Param		  status query string false "status"
//		@Success      200  {array}  models.Task
//		@Failure      400,404,500  {string}  stirng
//		@Router       /search-with-filters [get]
func (a *API) handleFindTask(c echo.Context) error {
	p := models.Task{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Assignee:    c.FormValue("assignee"),
		Status:      c.FormValue("status"),
	}

	if p.Name == "" && p.Description == "" && p.Assignee == "" && p.Status == "" {
		return c.JSON(http.StatusBadRequest, atLeastOneParamRequired)
	}

	results, err := a.r.GetMatchingTasks(c.Request().Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(results) == 0 {
		return c.JSON(http.StatusNotFound, noTasksFound)
	}

	return c.JSON(http.StatusOK, results)
}

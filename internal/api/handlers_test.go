package api

import (
	"context"
	"testing"

	"github.com/Marsredskies/todo-list/internal/db"
	"github.com/Marsredskies/todo-list/internal/envconfig"
	"github.com/Marsredskies/todo-list/internal/models"
	"github.com/stretchr/testify/require"
)

func TestSaveAndUpdateTask(t *testing.T) {
	ctx := context.Background()

	cnf := envconfig.MustGetConfig()

	dbConn, err := db.ConnectDB(ctx, cnf)
	require.NoError(t, err)

	db.DropMigrations(dbConn)
	db.MustApplyMigrations(ctx, cnf)

	createParams := models.Task{
		Name:        "test task",
		Description: "description for the test task",
		Assignee:    "smbd",
		Status:      "",
	}

	require.NoError(t, createParams.Validate())

	api := MustInitNewAPI(ctx, cnf)

	id, err := api.saveTaskToDb(ctx, createParams)
	require.NoError(t, err)

	require.Equal(t, int64(1), id)

	updateParams := models.Task{
		ID:          id,
		Name:        "test task update",
		Description: "updated description for the test task",
		Assignee:    "smbd else",
		Status:      "done",
	}

	err = api.updateTaskById(ctx, updateParams)
	require.NoError(t, err)
}

func TestFindTasks(t *testing.T) {
	ctx := context.Background()

	cnf := envconfig.MustGetConfig()

	dbConn, err := db.ConnectDB(ctx, cnf)
	require.NoError(t, err)

	db.DropMigrations(dbConn)
	db.MustApplyMigrations(ctx, cnf)

	createParams1 := models.Task{
		Name:        "first task",
		Description: "description description task",
		Assignee:    "person 1",
		Status:      "done",
	}

	createParams2 := models.Task{
		Name:        "secont task",
		Description: "description for the second task",
		Assignee:    "person 2",
		Status:      "done",
	}

	require.NoError(t, createParams1.Validate())
	require.NoError(t, createParams2.Validate())

	api := MustInitNewAPI(ctx, cnf)

	_, err = api.saveTaskToDb(ctx, createParams1)
	require.NoError(t, err)

	_, err = api.saveTaskToDb(ctx, createParams2)
	require.NoError(t, err)

	searchParams := models.Task{
		Name:        "task",
		Description: "description",
		Assignee:    "person",
	}

	results, err := api.getMatchingTasks(ctx, searchParams)
	require.NoError(t, err)
	require.Len(t, results, 2)
}

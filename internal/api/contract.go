package api

import (
	"context"

	"github.com/Marsredskies/todo-list/internal/models"
)

type Repository interface {
	GetMatchingTasks(ctx context.Context, params models.Task) ([]models.Task, error)
	SaveTask(ctx context.Context, params models.Task) (int64, error)
	UpdateTaskById(ctx context.Context, params models.Task) error
	DeleteById(ctx context.Context, id int64) error
	CheckIfTaskExists(id int64) (bool, error)
}

package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Marsredskies/todo-list/internal/models"
)

type TasksRepo struct {
	db *DB
}

func NewTaskRepo(db *DB) *TasksRepo {
	return &TasksRepo{db}
}

func (t *TasksRepo) GetMatchingTasks(ctx context.Context, params models.Task) ([]models.Task, error) {
	var results []models.Task
	query, args, err := params.SqlSelectLike()
	if err != nil {
		return nil, err
	}

	err = t.db.SelectCtx(ctx, &results, query, args...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (t *TasksRepo) SaveTask(ctx context.Context, params models.Task) (int64, error) {
	query, args, err := params.SqlInsert()
	if err != nil {
		return 0, err
	}

	return t.db.ExecCtxReturningId(ctx, query, args...)
}

func (t *TasksRepo) UpdateTaskById(ctx context.Context, params models.Task) error {
	query, args, err := params.SqlUpdate()
	if err != nil {
		return err
	}

	err = t.db.ExecCtx(ctx, query, args...)
	return err
}

func (t *TasksRepo) CheckIfTaskExists(id int64) (bool, error) {
	var task models.Task
	err := t.db.GetById(&task,
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

func (t *TasksRepo) DeleteById(ctx context.Context, id int64) error {
	err := t.db.ExecCtx(ctx,
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

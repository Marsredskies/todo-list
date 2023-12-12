package models

import "errors"

type Task struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type Request struct {
	Name        string `json:"name"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

var (
	errEmptyTaskName        = errors.New("task name can't be empty")
	errEmptyTaskDescription = errors.New("task description can't be empty")
	errInvalidTaskStatus    = errors.New(`invalid task status, can be either "to do" or "in progress" or "done"`)
)

func (t *Task) Validate() error {
	if t.Name == "" {
		return errEmptyTaskName
	}

	if t.Description == "" {
		return errEmptyTaskDescription
	}

	return t.ValidateStatus()
}

func (t *Task) ValidateStatus() error {
	if t.Status != "" {
		if t.Status != todo && t.Status != inProgress && t.Status != done {
			return errInvalidTaskStatus
		}
		return nil
	} else {
		t.Status = todo
		return nil
	}
}

package models

import "errors"

type Task struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Assignee    string `json:"assignee" db:"assignee"`
	Status      string `json:"status" db:"status"`
	Description string `json:"description" db:"description"`
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
	errNoAssignee           = errors.New("task should be assigned to someone")
)

func (t *Task) Validate() error {
	if t.Name == "" {
		return errEmptyTaskName
	}

	if t.Description == "" {
		return errEmptyTaskDescription
	}

	if t.Assignee == "" {
		return errNoAssignee
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

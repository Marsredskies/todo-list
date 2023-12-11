package api

import "errors"

type Task struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return errors.New("task name can't be empty")
	}

	if t.Description == "" {
		return errors.New("task description can't be empty")
	}

	if t.Status != "" {
		if t.Status != todo && t.Status != inProgress && t.Status != done {
			return errors.New(`invalid task status, can be either "to do" or "in progress" or "done"`)
		}
	}

	return nil
}

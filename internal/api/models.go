package api

import "errors"

type Params struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

func (p *Params) Validate() error {
	if p.Name == "" {
		return errors.New("task name can't be empty")
	}

	if p.Description == "" {
		return errors.New("task description can't be empty")
	}

	if p.Status != "" {
		if p.Status != todo && p.Status != inProgress && p.Status != done {
			return errors.New(`invalid task status, can be either "to do" or "in progress" or "done"`)
		}
	}

	return nil
}

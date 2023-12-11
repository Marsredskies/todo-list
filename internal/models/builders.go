package models

import (
	sq "github.com/Masterminds/squirrel"
)

func (t *Task) SqlInsert() (string, []interface{}, error) {
	values := sq.Eq{
		"name":        t.Name,
		"description": t.Description,
		"status":      todo,
	}

	if t.Assignee != "" {
		values["assignee"] = t.Assignee
	}

	return sq.Insert("public.tasks").PlaceholderFormat(sq.Dollar).SetMap(values).Suffix("RETURNING id").ToSql()
}

func (t *Task) SqlUpdate() (string, []interface{}, error) {
	b := sq.Update("public.tasks")
	b = t.Set(b)
	b = b.Where("id = ?", t.ID)

	return b.PlaceholderFormat(sq.Dollar).ToSql()
}

func (t *Task) Set(b sq.UpdateBuilder) sq.UpdateBuilder {
	if t.Name != "" {
		b = b.Set("name", t.Name)
	}

	if t.Description != "" {
		b = b.Set("description", t.Description)
	}

	if t.Assignee != "" {
		b = b.Set("assignee", t.Assignee)
	}

	if t.Status != "" {
		b = b.Set("status", t.Status)
	}

	b = b.Set("updated_at", sq.Expr("now()"))

	return b
}

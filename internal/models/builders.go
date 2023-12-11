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

func (t *Task) SqlSelectLike() (string, []interface{}, error) {
	b := sq.Select("id", "name", "description", "status", "assignee").From("public.tasks")

	if t.Name != "" {
		b = b.Where("name LIKE ?", likeParam(t.Name))
	}

	if t.Description != "" {
		b = b.Where("description LIKE ?", likeParam(t.Description))
	}

	if t.Assignee != "" {
		b = b.Where("assignee LIKE ?", likeParam(t.Assignee))
	}

	if t.Status != "" {
		b = b.Where("status LIKE ?", likeParam(t.Status))
	}

	b = b.Where("deleted_at IS NULL")
	return b.PlaceholderFormat(sq.Dollar).ToSql()
}

func likeParam(v string) string {
	return "%" + v + "%"
}

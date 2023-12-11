package todolist

import "github.com/Marsredskies/todo-list/internal/db"

func init() {
	db.AddMigration(1, m1)
}

var m1 = `
	CREATE TABLE public.todo_list (
		id					BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY
		,name         		TEXT NOT NULL
		,description		TEXT
		,assignee			TEXT				
		,status				TEXT NOT NULL DEFAULT 'to do'
		,created_at			TIMESTAMPTZ NOT NULL DEFAULT now()
		,updated_at			TIMESTAMPTZ
		,deleted_at			TIMESTAMPTZ
	);
`
package db

import "context"

func (d *DB) ExecCtxReturningId(ctx context.Context, query string, args ...interface{}) (int64, error) {
	var id int64
	err := d.db.QueryRowxContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DB) ExecCtx(ctx context.Context, query string, args ...interface{}) error {
	_, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) SelectCtx(ctx context.Context, dest any, query string, args ...interface{}) error {
	err := d.db.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return err
	}

	return nil
}
func (d *DB) GetById(dest any, query string, args ...interface{}) error {
	err := d.db.Get(dest, query, args...)
	if err != nil {
		return err
	}

	return nil
}

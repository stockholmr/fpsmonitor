package model

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB

	InstallCmd string
	CreateCmd  string
	UpdateCmd  string
	SelectCmd  string
	DeleteCmd  string
}

func (r *Repository) Install(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, r.InstallCmd)

	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Select(ctx context.Context, data interface{}, id interface{}) error {
	stmt, err := r.db.PreparexContext(ctx, r.SelectCmd)

	if err != nil {
		return err
	}

	err = stmt.GetContext(
		ctx,
		&data,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, data ...interface{}) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return -1, err
	}

	stmt, err := tx.PreparexContext(ctx, r.CreateCmd)

	if err != nil {
		return -1, err
	}

	result, err := stmt.ExecContext(ctx, data...)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	tx.Commit()
	id, _ := result.LastInsertId()
	return id, nil
}

func (r *Repository) Update(ctx context.Context, data ...interface{}) error {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, r.UpdateCmd)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, data...)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) Delete(ctx context.Context, data ...interface{}) error {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(ctx, r.DeleteCmd)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, data...)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

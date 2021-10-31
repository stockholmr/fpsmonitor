package computer

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type Computer struct {
	ID      null.Int    `db:"id" json:"id"`
	Created null.String `db:"created" json:"created"`
	Updated null.String `db:"updated" json:"updated"`
	Deleted null.String `db:"deleted" json:"deleted"`

	Name null.String `db:"name" json:"name"`
}

type ComputerRepository interface {
	Install(context.Context) error
	Select(context.Context, string) (*Computer, error)
	Create(context.Context, *Computer) (int64, error)
	Update(context.Context, *Computer) error
	Delete(context.Context, int) error
	List(context.Context, int, int) ([]Computer, error)
}

type computerRepository struct {
	db *sqlx.DB
}

func NewComputerRepository(db *sqlx.DB) ComputerRepository {
	return &computerRepository{
		db: db,
	}
}

func (r *computerRepository) Install(ctx context.Context) error {
	_, err := r.db.ExecContext(
		ctx,
		`CREATE TABLE computers (
            "id" INTEGER,
            "created" TEXT,
            "updated" TEXT,
            "deleted" TEXT,
            "name" TEXT,
            PRIMARY KEY("id" AUTOINCREMENT)
        )`,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *computerRepository) Select(ctx context.Context, id string) (*Computer, error) {
	data := Computer{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT
            id,
            created,
            updated,
            deleted,
            name
        FROM computers
        WHERE name=?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.GetContext(
		ctx,
		&data,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *computerRepository) Create(ctx context.Context, data *Computer) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return -1, err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`INSERT INTO computers (
            created,
            name
        ) VALUES (?,?)`,
	)

	if err != nil {
		return -1, err
	}

	result, err := stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		data.Name,
	)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	tx.Commit()
	id, _ := result.LastInsertId()
	return id, nil
}

func (r *computerRepository) Update(ctx context.Context, data *Computer) error {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computers SET
            updated=?,
            name=?
        WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		data.Name,
		data.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *computerRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computers SET
            deleted=?
        WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		id,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *computerRepository) List(ctx context.Context, start int, count int) ([]Computer, error) {
	data := []Computer{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT
            id,
            created,
            updated,
            deleted,
            name
        FROM computers
        LIMIT ?, ?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.SelectContext(
		ctx,
		&data,
		start,
		count,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

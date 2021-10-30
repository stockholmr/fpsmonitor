package computer

import (
	"context"
	"database/sql"
	"fpsmonitor/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type ComputerRepository interface {
	Install(ctx context.Context) error
	Get(ctx context.Context, computername string) (*Computer, error)
	Create(ctx context.Context, data *Computer) (int64, error)
	Update(ctx context.Context, data *Computer) error
	Delete(ctx context.Context, ID int) error
}

func NewComputerRepository(db *sqlx.DB) ComputerRepository {
	return &model.Repository{
		db: db,

		InstallCmd: ,

	}
}

func (c *computerRepository) Install(ctx context.Context) error {
	_, err := c.db.ExecContext(
		ctx,
		`CREATE TABLE "computers" (
            "id" INTEGER,
            "name" TEXT,
            "created" TEXT,
			"updated" TEXT,
			"deleted" TEXT,
            PRIMARY KEY("id" AUTOINCREMENT)
        )`,
	)

	if err != nil {
		return err
	}
	return nil
}

func (c *computerRepository) Get(ctx context.Context, computername string) (*Computer, error) {
	data := Computer{}

	stmt, err := c.db.PreparexContext(
		ctx,
		`SELECT 
            id,
            name,
            created,
			updated,
			deleted
        FROM computers 
        WHERE name=?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.GetContext(
		ctx,
		&data,
		computername,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (c *computerRepository) Create(ctx context.Context, data *Computer) (int64, error) {
	tx, err := c.db.BeginTxx(ctx, nil)

	if err != nil {
		return -1, err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`INSERT INTO computers(
            name,
			created
        ) VALUES (?,?)`,
	)

	if err != nil {
		return -1, err
	}

	result, err := stmt.ExecContext(
		ctx,
		data.Name,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	tx.Commit()
	id, _ := result.LastInsertId()
	return id, nil
}

func (c *computerRepository) Update(ctx context.Context, data *Computer) error {
	tx, err := c.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computers SET
            name=?,
			updated=?
        WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		data.Name,
		time.Now().Format("2006-01-02 15:04:05"),
		data.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *computerRepository) Delete(ctx context.Context, ID int) error {
	tx, err := c.db.BeginTxx(ctx, nil)

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
		ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

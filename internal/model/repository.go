package computer

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type {Model}Repository interface {
	Install(ctx context.Context) error
	Get(ctx context.Context, computerID int) (*[]{Model}, error)
	Create(ctx context.Context, data *{Model}) (int64, error)
	Update(ctx context.Context, data *{Model}) error
	Delete(ctx context.Context, ID int) error
}

type {model}Repository struct {
	db *sqlx.DB
}

func New{Model}Repository(db *sqlx.DB) {Model}Repository {
	return &{model}Repository{
		db: db,
	}
}

func (n *networkAdapterRepository) Install(ctx context.Context) error {
	_, err := n.db.ExecContext(
		ctx,
		`CREATE TABLE "computer_network_adapters" (
            "id" INTEGER,
            "computer_id" INTEGER,
			"name" TEXT,
            "mac_address" TEXT,
			"ip_address" TEXT,
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

func (n *networkAdapterRepository) Get(ctx context.Context, computerID int) (*[]NetworkAdapter, error) {
	data := []NetworkAdapter{}

	stmt, err := n.db.PreparexContext(
		ctx,
		`SELECT 
            id,
            name,
			mac_address,
            ip_address,
            created,
			updated,
			deleted
        FROM computer_network_adapters 
        WHERE computer_id=?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.SelectContext(
		ctx,
		&data,
		computerID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (n *networkAdapterRepository) Create(ctx context.Context, data *NetworkAdapter) (int64, error) {
	tx, err := n.db.BeginTxx(ctx, nil)

	if err != nil {
		return -1, err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`INSERT INTO computer_network_adapters(
			computer_id,
            name,
			mac_address,
            ip_address,
			created
        ) VALUES (?,?,?,?,?)`,
	)

	if err != nil {
		return -1, err
	}

	result, err := stmt.ExecContext(
		ctx,
		data.ComputerID,
		data.Name,
		data.MacAddress,
		data.IPAddress,
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

func (n *networkAdapterRepository) Update(ctx context.Context, data *NetworkAdapter) error {
	tx, err := n.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computer_network_adapters SET
            name=?,
            ip_address=?,
			updated=?
        WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		data.Name,
		data.IPAddress,
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

func (n *networkAdapterRepository) Delete(ctx context.Context, ID int) error {
	tx, err := n.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computer_network_adapters SET
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

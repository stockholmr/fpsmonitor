package computer

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type NetworkAdapter struct {
	ID      null.Int    `db:"id" json:"id"`
	Created null.String `db:"created" json:"created"`
	Updated null.String `db:"updated" json:"updated"`
	Deleted null.String `db:"deleted" json:"deleted"`

	ComputerID null.Int    `db:"computer_id" json:"-"`
	Name       null.String `db:"name" json:"name"`
	MacAddress null.String `db:"mac_address" json:"mac_address"`
	IPAddress  null.String `db:"ip_address" json:"ip_address"`
}

type NetworkAdapterRepository interface {
	Install(context.Context) error
	Select(context.Context, int) (*NetworkAdapter, error)
	Create(context.Context, *NetworkAdapter) (int64, error)
	Update(context.Context, *NetworkAdapter) error
	Delete(context.Context, int) error
	List(context.Context, int, int) ([]NetworkAdapter, error)
}

type networkAdapterRepository struct {
	db *sqlx.DB
}

func NewNetworkAdapterRepository(db *sqlx.DB) NetworkAdapterRepository {
	return &networkAdapterRepository{
		db: db,
	}
}

func (r *networkAdapterRepository) Install(ctx context.Context) error {
	_, err := r.db.ExecContext(
		ctx,
		`CREATE TABLE computer_network_adapters (
            "id" INTEGER,
            "created" TEXT,
            "updated" TEXT,
            "deleted" TEXT,
            "computer_id" INTEGER NOT NULL,
			"name" TEXT,
            "mac_address" TEXT,
			"ip_address" TEXT,
			FOREIGN KEY("computer_id") REFERENCES "computers"("id") ON DELETE CASCADE ON UPDATE NO ACTION,
            PRIMARY KEY("id" AUTOINCREMENT)
        )`,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *networkAdapterRepository) Select(ctx context.Context, id int) (*NetworkAdapter, error) {
	data := NetworkAdapter{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT 
            id,
            created,
            updated,
            deleted,
            name,
			mac_address,
            ip_address
        FROM computer_network_adapters
        WHERE id=?`,
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

func (r *networkAdapterRepository) Create(ctx context.Context, data *NetworkAdapter) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return -1, err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`INSERT INTO computer_network_adapters (
            created,
            computer_id,
            name,
			mac_address,
            ip_address
        ) VALUES (?,?,?,?,?)`,
	)

	if err != nil {
		return -1, err
	}

	result, err := stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		data.ComputerID,
		data.Name,
		data.MacAddress,
		data.IPAddress,
	)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	tx.Commit()
	id, _ := result.LastInsertId()
	return id, nil
}

func (r *networkAdapterRepository) Update(ctx context.Context, data *NetworkAdapter) error {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computer_network_adapters SET
            updated=?,
            name=?,
            ip_address=?
        WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		data.Name,
		data.IPAddress,
		data.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *networkAdapterRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTxx(ctx, nil)

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
		id,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *networkAdapterRepository) List(ctx context.Context, start int, count int) ([]NetworkAdapter, error) {
	data := []NetworkAdapter{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT
            id,
            created,
            updated,
            deleted,
            name,
            mac_address,
            ip_address
        FROM computer_network_adapters
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

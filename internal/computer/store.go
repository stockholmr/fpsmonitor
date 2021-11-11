package computer

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type computerStore struct {
	db *sqlx.DB
}

type ComputerStoreInterface interface {
	Create(ctx context.Context, m *ComputerModel) (null.Int, error)
	UpdateLastActivityAt(ctx context.Context, m *ComputerModel) error
	Get(ctx context.Context, id string) (*ComputerModel, error)
	GetAll(ctx context.Context, start int, count int) ([]ComputerModel, error)
	SoftDelete(ctx context.Context, m *ComputerModel) error
	HardDelete(ctx context.Context, m *ComputerModel) error
}

func NewComputerStore(db *sqlx.DB) ComputerStoreInterface {
	return &computerStore{
		db: db,
	}
}

func (s *computerStore) Create(ctx context.Context, m *ComputerModel) (null.Int, error) {
	stmt, err := s.db.PrepareContext(
		ctx,
		`INSERT INTO computers (
            created,
            name
        ) VALUES (?,?)`,
	)

	if err != nil {
		return null.Int{}, err
	}

	result, err := stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		m.Name,
	)

	if err != nil {
		return null.Int{}, err
	}

	id, _ := result.LastInsertId()
	return null.IntFrom(id), nil
}

func (s *computerStore) UpdateLastActivityAt(ctx context.Context, m *ComputerModel) error {
	stmt, err := s.db.PrepareContext(
		ctx,
		`UPDATE computers SET
			updated=?
		WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		m.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *computerStore) Get(ctx context.Context, id string) (*ComputerModel, error) {
	model := ComputerModel{}

	stmt, err := s.db.PreparexContext(
		ctx,
		`SELECT * FROM computers WHERE name=?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.GetContext(
		ctx,
		&model,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (s *computerStore) GetAll(ctx context.Context, start int, count int) ([]ComputerModel, error) {
	model := make([]ComputerModel, 0)

	stmt, err := s.db.PreparexContext(
		ctx,
		`SELECT * FROM computers LIMIT ? OFFSET ?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.GetContext(
		ctx,
		&model,
		count,
		start,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return model, nil
}

func (s *computerStore) SoftDelete(ctx context.Context, m *ComputerModel) error {
	stmt, err := s.db.PreparexContext(
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
		m.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *computerStore) HardDelete(ctx context.Context, m *ComputerModel) error {
	stmt, err := s.db.PreparexContext(
		ctx,
		`DELETE FROM computers WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		m.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

type networkAdapterStore struct {
	db *sqlx.DB
}

type NetworkAdapterInterface interface {
	Create(ctx context.Context, m *NetworkAdapterModel) (null.Int, error)
	UpdateIPAddress(ctx context.Context, m *NetworkAdapterModel) error
	GetAllByComputerID(ctx context.Context, id int) ([]NetworkAdapterModel, error)
	SoftDelete(ctx context.Context, m *NetworkAdapterModel) error
	HardDelete(ctx context.Context, m *NetworkAdapterModel) error
}

func NewNetworkAdapterStore(db *sqlx.DB) ComputerStoreInterface {
	return &computerStore{
		db: db,
	}
}

func (s *networkAdapterStore) Create(ctx context.Context, m *NetworkAdapterModel) (null.Int, error) {
	stmt, err := s.db.PrepareContext(
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
		return null.Int{}, err
	}

	result, err := stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		m.ComputerID,
		m.Name,
		m.MacAddress,
		m.IPAddress,
	)

	if err != nil {
		return null.Int{}, err
	}

	id, _ := result.LastInsertId()
	return null.IntFrom(id), nil
}

func (s *networkAdapterStore) Update(ctx context.Context, m *NetworkAdapterModel) error {
	stmt, err := s.db.PrepareContext(
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
		m.Name,
		m.IPAddress,
		m.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type userStore struct {
	db *sqlx.DB
}

type UserStoreInterface interface {
	Create(ctx context.Context, m *UserModel) (null.Int, error)
	UpdateLastActivityAt(ctx context.Context, m *UserModel) error
	Get(ctx context.Context, id null.Int) (*UserModel, error)
	GetByUsername(ctx context.Context, id string) (*UserModel, error)
	SoftDelete(ctx context.Context, m *UserModel) error
	HardDelete(ctx context.Context, m *UserModel) error
}

func NewUserStore(db *sqlx.DB) UserStoreInterface {
	return &userStore{
		db: db,
	}
}

func (s *userStore) Create(ctx context.Context, m *UserModel) (null.Int, error) {
	stmt, err := s.db.PrepareContext(
		ctx,
		`INSERT INTO computer_users (
            created,
            computer_id,
            username
        ) VALUES (?,?,?)`,
	)

	if err != nil {
		return null.Int{}, err
	}

	result, err := stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		m.ID,
		m.Username,
	)

	if err != nil {
		return null.Int{}, err
	}

	id, _ := result.LastInsertId()
	return null.IntFrom(id), nil
}

func (s *userStore) UpdateLastActivityAt(ctx context.Context, m *UserModel) error {
	stmt, err := s.db.PrepareContext(
		ctx,
		`UPDATE users SET
			updated=?,
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

func (s *userStore) Get(ctx context.Context, id null.Int) (*UserModel, error) {
	model := UserModel{}

	stmt, err := s.db.PreparexContext(
		ctx,
		`SELECT * FROM users 
			WHERE id=?`,
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

func (s *userStore) GetByUsername(ctx context.Context, id string) (*UserModel, error) {
	model := UserModel{}

	stmt, err := s.db.PreparexContext(
		ctx,
		`SELECT * FROM users 
			WHERE username=?`,
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

func (s *userStore) SoftDelete(ctx context.Context, m *UserModel) error {
	stmt, err := s.db.PreparexContext(
		ctx,
		`UPDATE computer_users SET
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

func (s *userStore) HardDelete(ctx context.Context, m *UserModel) error {
	stmt, err := s.db.PreparexContext(
		ctx,
		`DELETE FROM computer_users WHERE id=?`,
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

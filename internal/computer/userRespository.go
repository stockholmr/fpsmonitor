package computer

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Install(ctx context.Context) error
	Get(ctx context.Context, username string) (*[]User, error)
	GetWithComputerID(ctx context.Context, username string, computerID int) (*User, error)
	Create(ctx context.Context, data *User) (int64, error)
	Update(ctx context.Context, data *User) error
	Delete(ctx context.Context, ID int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Install(ctx context.Context) error {
	_, err := u.db.ExecContext(
		ctx,
		`CREATE TABLE "computer_users" (
            "id" INTEGER,
            "computer_id" INTEGER,
			"username" TEXT,
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

func (u *userRepository) Get(ctx context.Context, username string) (*[]User, error) {
	data := []User{}

	stmt, err := u.db.PreparexContext(
		ctx,
		`SELECT 
            id,
            computer_id,
			username,
            created,
			updated,
			deleted
        FROM computers 
        WHERE username=?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.SelectContext(
		ctx,
		&data,
		username,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (u *userRepository) GetWithComputerID(ctx context.Context, username string, computerID int) (*User, error) {
	data := User{}

	stmt, err := u.db.PreparexContext(
		ctx,
		`SELECT 
            id,
			username,
            created,
			updated,
			deleted
        FROM computer_users 
        WHERE username=? AND computer_id=?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.GetContext(
		ctx,
		&data,
		username,
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

func (u *userRepository) Create(ctx context.Context, data *User) (int64, error) {
	tx, err := u.db.BeginTxx(ctx, nil)

	if err != nil {
		return -1, err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`INSERT INTO computer_users (
			computer_id,
            username,
			created
        ) VALUES (?,?,?)`,
	)

	if err != nil {
		return -1, err
	}

	result, err := stmt.ExecContext(
		ctx,
		data.ComputerID,
		data.Username,
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

func (u *userRepository) Update(ctx context.Context, data *User) error {
	tx, err := u.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computer_users SET
			updated=?
        WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
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

func (u *userRepository) Delete(ctx context.Context, ID int) error {
	tx, err := u.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
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
		ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

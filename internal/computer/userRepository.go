package computer

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type User struct {
	ID         null.Int    `db:"id" json:"id"`
	Created    null.String `db:"created" json:"created"`
	Updated    null.String `db:"updated" json:"updated"`
	Deleted    null.String `db:"deleted" json:"deleted"`
	ComputerID null.Int    `db:"computer_id" json:"computer_id"`
	Username   null.String `db:"username" json:"username"`
}

type UserRepository interface {
	Install(context.Context) error
	Select(context.Context, int) (*User, error)
	Create(context.Context, *User) (int64, error)
	Update(context.Context, *User) error
	Delete(context.Context, int) error
	List(context.Context, int, int) ([]User, error)

	SelectWithUsername(context.Context, string) (*User, error)
	SelectWithUsernameAndComputerID(context.Context, int, string) (*User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Install(ctx context.Context) error {
	_, err := r.db.ExecContext(
		ctx,
		`CREATE TABLE computer_users (
            "id" INTEGER,
            "created" TEXT,
            "updated" TEXT,
            "deleted" TEXT,
            "computer_id" INTEGER NOT NULL,
			"username" TEXT,
			FOREIGN KEY("computer_id") REFERENCES "computers"("id") ON DELETE CASCADE ON UPDATE NO ACTION,
            PRIMARY KEY("id" AUTOINCREMENT)
        )`,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Select(ctx context.Context, id int) (*User, error) {
	data := User{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT 
            id,
            created,
            updated,
            deleted,
            computer_id,
            username
        FROM computer_users
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

func (r *userRepository) SelectWithUsername(ctx context.Context, id string) (*User, error) {
	data := User{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT 
            id,
            created,
            updated,
            deleted,
            computer_id,
            username
        FROM computer_users
        WHERE username=?`,
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

func (r *userRepository) SelectWithUsernameAndComputerID(ctx context.Context, id int, username string) (*User, error) {
	data := User{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT 
            id,
            created,
            updated,
            deleted,
            computer_id,
            username
        FROM computer_users
        WHERE computer_id=? AND username=?`,
	)

	if err != nil {
		return nil, err
	}

	err = stmt.GetContext(
		ctx,
		&data,
		id,
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

func (r *userRepository) Create(ctx context.Context, data *User) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return -1, err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`INSERT INTO computer_users (
            created,
            computer_id,
            username
        ) VALUES (?,?,?)`,
	)

	if err != nil {
		return -1, err
	}

	result, err := stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		data.ComputerID,
		data.Username,
	)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	tx.Commit()
	id, _ := result.LastInsertId()
	return id, nil
}

func (r *userRepository) Update(ctx context.Context, data *User) error {
	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PreparexContext(
		ctx,
		`UPDATE computer_users SET
            updated=?,
			username=?
        WHERE id=?`,
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		time.Now().Format("2006-01-02 15:04:05"),
		data.Username,
		data.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTxx(ctx, nil)

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
		id,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *userRepository) List(ctx context.Context, start int, count int) ([]User, error) {
	data := []User{}

	stmt, err := r.db.PreparexContext(
		ctx,
		`SELECT
            id,
            created,
            updated,
            deleted,
            computer_id,
            username
        FROM computer_users
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

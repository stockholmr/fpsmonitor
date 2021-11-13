package auth

import "github.com/jmoiron/sqlx"

func Migrate(db *sqlx.DB) error {
	migrations := []string{

		`CREATE TABLE users (
            id INTEGER,
            created TEXT,
            updated TEXT,
            deleted TEXT,
            username TEXT,
            password TEXT,
            PRIMARY KEY( id  AUTOINCREMENT)
        )`,
	}

	for _, s := range migrations {
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}

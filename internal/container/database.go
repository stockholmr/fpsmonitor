package container

import (
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func (c *container) InitDB(file string) {

	dir := path.Dir(file)
	err := os.Mkdir(dir, 0776)
	if err != nil {
		if !os.IsExist(err) {
			c.Fatal(err)
		}
	}

	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		c.Fatal(err)
	}

	c.db = db
}

func (c *container) SetDB(db *sqlx.DB) {
	c.db = db
}

func (c *container) DB() *sqlx.DB {
	if c.db == nil {
		db, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			c.Error(err)
		}
		c.db = db
	}

	return c.db
}

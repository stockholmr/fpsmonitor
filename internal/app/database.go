package app

import (
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func (a *app) InitDB(file string) {

	dir := path.Dir(file)
	err := os.Mkdir(dir, 0776)
	if err != nil {
		if !os.IsExist(err) {
			a.Fatal(err)
		}
	}

	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		a.Fatal(err)
	}

	a.db = db
}

func (a *app) SetDB(db *sqlx.DB) {
	a.db = db
}

func (a *app) DB() *sqlx.DB {
	if a.db == nil {
		db, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			a.Error(err)
		}
		a.db = db
	}

	return a.db
}

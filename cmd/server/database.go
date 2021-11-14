package main

import (
	"log"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stockholmr/fpsmonitor/internal/auth"
	"github.com/stockholmr/fpsmonitor/internal/computer"
)

func InitDB(file string, logg *log.Logger) *sqlx.DB {

	migrate := false

	dir := path.Dir(file)
	err := os.Mkdir(dir, 0776)
	if err != nil {
		if !os.IsExist(err) {
			logg.Fatal(err)
		}
	}

	_, err = os.Stat(file)
	if err != nil {
		if !os.IsExist(err) {
			migrate = true
		}
	}

	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		logg.Fatal(err)
	}

	if migrate {
		err = auth.Migrate(db)
		if err != nil {
			logg.Fatal(err)
		}

		err = computer.Migrate(db)
		if err != nil {
			logg.Fatal(err)
		}
	}

	return db
}

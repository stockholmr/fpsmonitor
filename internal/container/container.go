package container

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type container struct {
	log *log.Logger
	db  *sqlx.DB
}

type Container interface {
	Fatal(...interface{})
	Error(...interface{})
	Warning(...interface{})
	Info(...interface{})
	Debug(...interface{})

	SetLog(*log.Logger)
	InitFileLog(string)
	Log() *log.Logger

	SetDB(*sqlx.DB)
	InitDB(string)
	DB() *sqlx.DB
}

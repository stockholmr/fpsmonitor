package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type app struct {
	log          *log.Logger
	db           *sqlx.DB
	config       *ConfigModel
	sessionStore sessions.Store
	sessionKeys  *SessionKeys
	router       *mux.Router

	csrfMiddleware func(http.Handler) http.Handler
}

type App interface {
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

	InitConfig(string)
	SetConfig(*ConfigModel)
	Config() *ConfigModel

	InitSessionKeysFromBase64(string, string)
	InitSessionKeys([]byte, []byte)
	SetSessionKeys(*SessionKeys)
	SessionKeys() *SessionKeys
	SetSessionStore(sessions.Store)
	SessionStore() sessions.Store

	InitRouter()
	SetRouter(*mux.Router)
	Router() *mux.Router

	CsrfMiddleware(next http.Handler) http.Handler
}

func New() App {
	return &app{}
}

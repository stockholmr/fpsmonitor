package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type App struct {
	log            *log.Logger
	db             *sqlx.DB
	config         *ConfigModel
	sessionStore   sessions.Store
	sessionKeys    *SessionKeys
	router         *mux.Router
	controllers    map[string]interface{}
	csrfMiddleware func(http.Handler) http.Handler
	server         *http.Server
}

func New() *App {
	return &App{
		controllers: make(map[string]interface{}),
	}
}

package admin

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/stockholmr/fpsmonitor/internal/auth"
)

type controller struct {
	db   *sqlx.DB
	logg *log.Logger
	auth auth.Controller
}

type Controller interface {
	Index(http.ResponseWriter, *http.Request)
}

func Init(r *mux.Router, db *sqlx.DB, auth auth.Controller) Controller {
	c := &controller{
		db:   db,
		logg: log.Default(),
		auth: auth,
	}

	c.initLog()
	c.register(r)
	return c
}

func InitWithLogger(r *mux.Router, db *sqlx.DB, auth auth.Controller, logger *log.Logger) Controller {
	c := &controller{
		db:   db,
		logg: logger,
		auth: auth,
	}
	c.initLog()
	c.register(r)
	return c
}

func (c *controller) register(router *mux.Router) {
	r := router.PathPrefix("/admin").Subrouter()
	r.Handle("", c.auth.AuthenticateSession(http.HandlerFunc(c.Index))).Methods("GET").Name("admin")
}

func (c *controller) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin Portal"))
}

package admin

import (
	"net/http"

	"github.com/stockholmr/fpsmonitor/internal/app"
	"github.com/stockholmr/fpsmonitor/internal/auth"
)

type controller struct {
	app app.App
}

type Controller interface {
	Index(http.ResponseWriter, *http.Request)
}

func Init(app app.App) Controller {
	c := &controller{app: app}

	auth, ok := app.Controller("auth").(auth.Controller)
	if !ok {
		c.app.Fatal("missing controller: auth")
	}

	r := c.app.Router().PathPrefix("/admin").Subrouter()
	r.Handle("", auth.AuthenticateSession(http.HandlerFunc(c.Index))).Methods("GET").Name("admin")

	return c
}

func (c *controller) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin Portal"))
}

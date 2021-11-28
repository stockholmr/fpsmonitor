package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) InitRouter() {
	a.router = mux.NewRouter().StrictSlash(true)

	a.router.HandleFunc("/favicon.ico", a.Favicon()).Methods("GET").Name("favicon")
	a.router.HandleFunc("/bootstrap", a.Bootstrap()).Methods("GET").Name("bootstrap")
	a.router.HandleFunc("/jquery", a.Jquery()).Methods("GET").Name("jquery")
	a.router.HandleFunc("/axios", a.Axios()).Methods("GET").Name("axios")
	a.router.HandleFunc("/vuejsdev", a.VuejsDev()).Methods("GET").Name("vuejs_dev")
	a.router.HandleFunc("/vuejs", a.Vuejs()).Methods("GET").Name("vuejs")
}

func (a *App) SetRouter(router *mux.Router) {
	a.router = router
}

func (a *App) Router() *mux.Router {
	if a.router == nil {
		a.InitRouter()
	}

	return a.router
}

func (a *App) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

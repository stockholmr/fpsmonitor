package app

import "github.com/gorilla/mux"

func (a *app) InitRouter() {
	a.router = mux.NewRouter().StrictSlash(true)
}

func (a *app) SetRouter(router *mux.Router) {
	a.router = router
}

func (a *app) Router() *mux.Router {
	if a.router == nil {
		a.router = mux.NewRouter().StrictSlash(true)
	}

	return a.router
}

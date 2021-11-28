package app

import (
	"fmt"
	"net/http"
)

func (a *App) ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Forbidden"))
}

func (a *App) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Error"))
}

func (a *App) Favicon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(favicon())))
		w.Write(favicon())
	}
}

func (a *App) Axios() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(axios())))
		w.Write(axios())
	}
}

func (a *App) Bootstrap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bootstrap())))
		w.Write(bootstrap())
	}
}

func (a *App) Jquery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jquery())))
		w.Write(jquery())
	}
}

func (a *App) VuejsDev() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(vuejs_dev())))
		w.Write(vuejs_dev())
	}
}

func (a *App) Vuejs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(vuejs())))
		w.Write(vuejs())
	}
}

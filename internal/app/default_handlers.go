package app

import "net/http"

func (a *app) ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Forbidden"))
}

func (a *app) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Error"))
}

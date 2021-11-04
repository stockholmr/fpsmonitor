package auth

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
)

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func Login(db *sqlx.DB) http.Handler {
	return alice.New(SessionMiddleware(db), AuthenticationMiddleware(db)).ThenFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}

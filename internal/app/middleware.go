package app

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (a *App) InitCsrfMiddleware() {
	a.csrfMiddleware = csrf.Protect(
		a.SessionKeys().EncryptionKey,
		csrf.RequestHeader("Authenticity-Token"),
		csrf.FieldName("authenticity_token"),
		csrf.ErrorHandler(http.HandlerFunc(a.ForbiddenHandler)),
	)
}

func (a *App) CsrfMiddleware(next http.Handler) http.Handler {
	if a.csrfMiddleware == nil {
		a.csrfMiddleware = csrf.Protect(
			a.SessionKeys().EncryptionKey,
			csrf.RequestHeader("Authenticity-Token"),
			csrf.FieldName("authenticity_token"),
			csrf.ErrorHandler(http.HandlerFunc(a.ForbiddenHandler)),
		)
	}

	return a.csrfMiddleware(next)
}

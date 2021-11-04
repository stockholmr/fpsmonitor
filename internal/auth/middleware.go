package auth

import (
	"context"
	"fpsmonitor/internal/logging"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
)

type Key string

const UserKey Key = "user"
const SessionKey Key = "session"

func SessionMiddleware(db *sqlx.DB) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "DASHBOARD-SESSION")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			request := r.Clone(context.WithValue(r.Context(), SessionKey, session))
			next.ServeHTTP(w, request)
		})
	}
}

func AuthenticationMiddleware(db *sqlx.DB) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRepo := NewUserRepository(db)

			session := r.Context().Value(SessionKey).(*sessions.Session)
			userid := session.Values["userid"]
			if userid == nil {
				logging.Warning("session does not exist")
				redirect(w, r, "/login")
				return
			}

			user, err := userRepo.Select(r.Context(), userid.(int))
			if err != nil {
				logging.Warning("user does not exist")
				redirect(w, r, "/logout")
				return
			}

			/*if !user.Active {
				logging.Warning("user is disabled; session deleted")
				redirect(w, r, "/logout")
				return
			}*/

			request := r.Clone(context.WithValue(r.Context(), UserKey, user.ID))
			next.ServeHTTP(w, request)
		})
	}
}

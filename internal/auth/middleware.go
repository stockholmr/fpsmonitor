package auth

import (
	"context"
	"net/http"
)

func (c *controller) AuthenticateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, _, err := c.validateSession(r)
		if err != nil {
			c.Error(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if userID.Valid {
			request := r.Clone(context.WithValue(r.Context(), "user_id", userID.Int64))
			next.ServeHTTP(w, request)
			return
		}

		c.Redirect(w, r, "/login")
	})
}

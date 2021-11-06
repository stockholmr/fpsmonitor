package computer

import (
	"net/http"
)

type Key string

const UserKey Key = "user"
const SessionKey Key = "session"

func (c *computerController) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.log.Trace("%s|%s|%s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

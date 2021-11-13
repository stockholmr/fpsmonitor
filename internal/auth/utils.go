package auth

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/sessions"
	"gopkg.in/guregu/null.v3"
)

func (c *controller) toInteger(v interface{}) null.Int {
	if v == nil {
		return null.Int{}
	}

	if "string" == reflect.TypeOf(v).Name() {
		integer, err := strconv.Atoi(v.(string))
		if err != nil {
			return null.Int{}
		}
		return null.IntFrom(int64(integer))
	}

	if "int" == reflect.TypeOf(v).Name() {
		integer, ok := v.(int)
		if !ok {
			return null.Int{}
		}
		return null.IntFrom(int64(integer))
	}

	return null.Int{}
}

func (c *controller) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (c *controller) validateSession(r *http.Request) (null.Int, *sessions.Session, error) {
	// retrieve session from store
	session, err := c.sessionStore.Get(r, "AUTH_SESSION")
	if err != nil {
		return null.Int{}, nil, err
	}

	userStore := NewUserStore(c.db)

	// convert session value into integer
	userID := c.toInteger(session.Values["userid"])

	if !userID.Valid {

		// failed to convert session value to integer
		return null.Int{}, session, nil

	}

	// retrive user from database
	user, err := userStore.Get(r.Context(), userID)
	if err != nil || user == nil {

		// user does not exist or there was an error
		// retrieving the user from the database

		return null.Int{}, session, err

	}

	/*if !user.Active {
		c.log.Warn("user is disabled; session deleted")
		redirect(w, r, "/logout")
		return
	}*/

	// user is valid and the session is valid
	return user.ID, session, nil
}

package auth

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/sessions"
	"gopkg.in/guregu/null.v3"
)

func (c *controller) ToInteger(v interface{}) null.Int {
	if v == nil {
		return null.Int{}
	}

	if reflect.TypeOf(v).Name() == "string" {
		integer, err := strconv.Atoi(v.(string))
		if err != nil {
			return null.Int{}
		}
		return null.IntFrom(int64(integer))
	}

	if reflect.TypeOf(v).Name() == "int64" {
		integer, ok := v.(int64)
		if !ok {
			return null.Int{}
		}
		return null.IntFrom(integer)
	}

	if reflect.TypeOf(v).Name() == "int" {
		integer, ok := v.(int)
		if !ok {
			return null.Int{}
		}
		return null.IntFrom(int64(integer))
	}

	return null.Int{}
}

func (c *controller) GetSession(r *http.Request) (*sessions.Session, error) {
	// retrieve session from store
	session, err := c.app.SessionStore().Get(r, "AUTH_SESSION")
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (c *controller) ValidateSession(r *http.Request) (null.Int, *sessions.Session, error) {
	// retrieve session from store
	session, err := c.GetSession(r)
	if err != nil {
		return null.Int{}, nil, err
	}

	userStore := NewUserStore(c.app.DB())

	// convert session value into integer
	userID := c.ToInteger(session.Values["userid"])

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

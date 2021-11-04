package auth

import "github.com/gorilla/sessions"

var store *sessions.CookieStore

func Init(encKey string) {
	store = sessions.NewCookieStore([]byte(encKey))
}

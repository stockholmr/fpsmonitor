package app

import (
	"encoding/base64"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type SessionKeys struct {
	AuthenticationKey []byte
	EncryptionKey     []byte
}

func (a *app) InitSessionKeysFromBase64(authKey string, encKey string) {
	authKeyBytes, err := base64.StdEncoding.DecodeString(authKey)
	if err != nil {
		a.Error(err)
	}

	encKeyBytes, err := base64.StdEncoding.DecodeString(encKey)
	if err != nil {
		a.Error(err)
	}

	a.InitSessionKeys(authKeyBytes, encKeyBytes)
}

func (a *app) InitSessionKeys(authKey []byte, encKey []byte) {
	a.sessionKeys = &SessionKeys{
		AuthenticationKey: authKey,
		EncryptionKey:     encKey,
	}
}

func (a *app) SetSessionKeys(keys *SessionKeys) {
	a.sessionKeys = keys
}

func (a *app) SetSessionStore(store sessions.Store) {
	a.sessionStore = store
}

func (a *app) SessionKeys() *SessionKeys {
	if a.sessionKeys == nil {
		a.sessionKeys = &SessionKeys{
			AuthenticationKey: securecookie.GenerateRandomKey(64),
			EncryptionKey:     securecookie.GenerateRandomKey(32),
		}
	}

	return a.sessionKeys
}

func (a *app) SessionStore() sessions.Store {
	if a.sessionStore == nil {
		a.sessionStore = sessions.NewCookieStore(
			a.SessionKeys().AuthenticationKey,
			a.SessionKeys().EncryptionKey,
		)
	}

	return a.sessionStore
}

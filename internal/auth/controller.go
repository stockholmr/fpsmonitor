package auth

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type controller struct {
	db           *sqlx.DB
	log          *log.Logger
	sessionStore sessions.Store
	templates    *Templates
}

type Controller interface {
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
}

func Init(r *mux.Router, db *sqlx.DB) Controller {
	c := &controller{
		db:  db,
		log: log.Default(),
		sessionStore: sessions.NewCookieStore(
			securecookie.GenerateRandomKey(64),
			securecookie.GenerateRandomKey(32),
		),
		templates: InitTemplates(),
	}

	c.initLog()
	c.register(r)

	return c
}

func InitWithLogger(r *mux.Router, db *sqlx.DB, logger *log.Logger) Controller {
	c := &controller{
		db:  db,
		log: logger,
		sessionStore: sessions.NewCookieStore(
			securecookie.GenerateRandomKey(64),
			securecookie.GenerateRandomKey(32),
		),
		templates: InitTemplates(),
	}

	c.register(r)
	return c
}

func (c *controller) register(router *mux.Router) {

	csrfMiddleware := csrf.Protect(
		securecookie.GenerateRandomKey(32),
		csrf.RequestHeader("Authenticity-Token"),
		csrf.FieldName("authenticity_token"),
		csrf.ErrorHandler(http.HandlerFunc(c.Forbidden)),
	)

	r := router.PathPrefix("/user").Subrouter()
	r.Use(csrfMiddleware)

	r.HandleFunc("/login", c.Login).Methods("GET", "POST").Name("login")
	r.HandleFunc("/logout", c.Logout).Methods("GET").Name("logout")
	r.HandleFunc("/register", c.Register).Methods("GET", "POST").Name("register")

}

func (c *controller) Login(w http.ResponseWriter, r *http.Request) {

	userID, session, err := c.ValidateSession(r)
	if err != nil {
		c.Error(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if userID.Valid {
		// user and session is valid redirect to admin portal
		c.Redirect(w, r, "/admin")
	}

	userStore := NewUserStore(c.db)

	if r.Method == "POST" {

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		remember := r.PostFormValue("remember")

		user, err := userStore.GetByUsername(r.Context(), username)
		if err != nil {
			c.Debug(err)
			c.templates.Login(w, TemplateData{
				"Title": "Login",
				"Error": "Invalid Username or Password",
			})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password)) != nil {
			c.templates.Login(w, TemplateData{
				"Title": "Login",
				"Error": "Invalid Username or Password",
			})
			return
		}

		session.Values["userid"] = user.ID.Int64
		session.Options.MaxAge = 0
		if remember != "" {
			// remember session for 1 week
			session.Options.MaxAge = 604800
		}

		err = session.Save(r, w)
		if err != nil {
			c.Warning("failed to save session")
		}

		userStore.UpdateLastActivityAt(r.Context(), user)
		c.Redirect(w, r, "/admin")

	}

	c.templates.Login(w, TemplateData{
		"Title": "Login",
	})
}

func (c *controller) Logout(w http.ResponseWriter, r *http.Request) {

	// Retrieve session from store
	session, err := c.sessionStore.Get(r, "AUTH_SESSION")
	if err != nil {
		c.Error(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	if session.Save(r, w) == nil {
		c.Redirect(w, r, "/login")
		return
	} else {
		c.Debug("failed to save session")
	}
}

func (c *controller) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		_ = ""
	}

	c.templates.Register(w, TemplateData{
		"Title": "Register",
	})
}

func (c *controller) Forbidden(w http.ResponseWriter, r *http.Request) {

}

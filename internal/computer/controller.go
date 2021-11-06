package computer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jcelliott/lumber"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	"gopkg.in/guregu/null.v3"
)

type computerController struct {
	log                lumber.Logger
	sessionStore       sessions.Store
	router             *mux.Router
	computerRepo       ComputerRepository
	networkAdapterRepo NetworkAdapterRepository
	userRepo           UserRepository
}

type ComputerController interface {
	Update(http.ResponseWriter, *http.Request)
}

func NewComputerController(db *sqlx.DB, log lumber.Logger, router *mux.Router, middleware ...alice.Constructor) ComputerController {
	c := &computerController{
		log:                log,
		router:             router,
		computerRepo:       NewComputerRepository(db),
		networkAdapterRepo: NewNetworkAdapterRepository(db),
		userRepo:           NewUserRepository(db),
	}

	m := []alice.Constructor{
		c.LoggingMiddleware,
	}
	m = append(m, middleware...)

	r := c.router.PathPrefix("/computers").Subrouter()
	r.Handle("/update", alice.New(m...).ThenFunc(c.Update)).Methods("POST").Name("update")
	r.Handle("/list", alice.New(m...).ThenFunc(c.Update)).Methods("POST").Name("list")
	r.Handle("/stylesheet", alice.New(m...).ThenFunc(c.Stylesheet)).Methods("GET")

	return c
}

func (c *computerController) Stylesheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(stylesheet())))
	w.Write(stylesheet())
}

func (c *computerController) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 10))
	defer cancel()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.log.Error("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var record struct {
		ID       null.Int
		Name     null.String      `json:"name"`
		Username null.String      `json:"username"`
		Adapters []NetworkAdapter `json:"adapters"`
	}

	err = json.Unmarshal(data, &record)
	if err != nil {
		c.log.Error("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var compID int64

	comp, err := c.computerRepo.Select(ctx, record.Name.String)
	if err != nil {
		c.log.Error("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comp != nil {

		// Computer record exists update the updated date field
		compID = comp.ID.Int64
		err := c.computerRepo.Update(ctx, comp)
		if err != nil {
			c.log.Error("%s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {

		// Create new computer record
		compID, err = c.computerRepo.Create(ctx, &Computer{
			Name: record.Name,
		})

		if err != nil {
			c.log.Error("%s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	// Create new user record
	_, err = c.userRepo.Create(ctx, &User{
		Username:   record.Username,
		ComputerID: null.IntFrom(compID),
	})

	if err != nil {
		c.log.Error("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	networkAdapters, err := c.networkAdapterRepo.SelectWithComputerID(ctx, int(compID))
	if err != nil {
		c.log.Error("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(networkAdapters) > 0 {

		for _, na := range networkAdapters {
			for _, nar := range record.Adapters {
				if na.MacAddress == nar.MacAddress {
					na.IPAddress = nar.IPAddress
					na.Name = nar.Name
					err = c.networkAdapterRepo.Update(ctx, &na)
					if err != nil {
						c.log.Error("%s", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}
		}

	} else {

		for _, na := range record.Adapters {
			na.ComputerID = null.IntFrom(compID)
			if _, err = c.networkAdapterRepo.Create(ctx, &na); err != nil {
				c.log.Error("%s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	}

	w.WriteHeader(http.StatusOK)
}

func (c *computerController) List(w http.ResponseWriter, r *http.Request) {

	list, err := c.userRepo.ListWithComputerNames(r.Context(), 0, 20)
	if err != nil {
		c.log.Error("%s", err)
		c.log.Trace("%s", err.(*ErrorEx).Func)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Records []User
	}{
		Title:   "Computer List",
		Records: list,
	}

	editorPage().ExecuteTemplate(w, "page", &data)
}

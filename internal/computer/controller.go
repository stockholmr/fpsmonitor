package computer

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type controller struct {
	db  *sqlx.DB
	logg *log.Logger
}

type Controller interface {
	Update(http.ResponseWriter, *http.Request)
}

func Init(r *mux.Router, db *sqlx.DB) Controller {
	c := &controller{
		db:  db,
		logg: log.Default(),
	}

	c.initLog()

	router := r.PathPrefix("/computers").Subrouter()

	router.HandleFunc("/update", c.Update).Methods("POST").Name("update")

	/*r.Handle("/update", alice.New(m...).ThenFunc(c.Update)).Methods("POST").Name("update")
	r.Handle("/list", alice.New(m...).ThenFunc(c.Update)).Methods("POST").Name("list")
	r.Handle("/stylesheet", alice.New(m...).ThenFunc(c.Stylesheet)).Methods("GET")*/

	return c
}

func InitWithLogger(r *mux.Router, db *sqlx.DB, logger *log.Logger) Controller {
	c := &controller{
		db:  db,
		logg: logger,
	}

	router := r.PathPrefix("/computers").Subrouter()

	router.HandleFunc("/update", c.Update).Methods("POST").Name("update")

	return c
}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 10))
	defer cancel()

	compStore := NewComputerStore(c.db)
	netStore := NewNetworkAdapterStore(c.db)
	userStore := NewUserStore(c.db)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var record struct {
		ID       null.Int
		Name     null.String           `json:"name"`
		Username null.String           `json:"username"`
		Adapters []NetworkAdapterModel `json:"adapters"`
	}

	err = json.Unmarshal(data, &record)
	if err != nil {
		c.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var compID null.Int

	comp, err := compStore.Get(ctx, record.Name.String)
	if err != nil {
		c.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comp != nil {

		// Computer record exists update the updated date field
		compID = null.IntFrom(comp.ID.Int64)
		err := compStore.UpdateLastActivityAt(ctx, comp)
		if err != nil {
			c.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {

		// Create new computer record
		compID, err = compStore.Create(ctx, &ComputerModel{
			Name: record.Name,
		})

		if err != nil {
			c.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	// Create new user record
	_, err = userStore.Create(ctx, &UserModel{
		Username:   record.Username,
		ComputerID: compID,
	})

	if err != nil {
		c.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	networkAdapters, err := netStore.GetAllByComputerID(ctx, int(compID.Int64))
	if err != nil {
		c.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(networkAdapters) > 0 {

		for _, na := range networkAdapters {
			for _, nar := range record.Adapters {
				if na.MacAddress == nar.MacAddress {
					na.IPAddress = nar.IPAddress
					na.Name = nar.Name
					err = netStore.Update(ctx, &na)
					if err != nil {
						c.Error(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}
		}

	} else {

		for _, na := range record.Adapters {
			na.ComputerID = compID
			if _, err = netStore.Create(ctx, &na); err != nil {
				c.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	}

	w.WriteHeader(http.StatusOK)
}

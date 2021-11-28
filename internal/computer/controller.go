package computer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/stockholmr/fpsmonitor/internal/app"
	"github.com/stockholmr/fpsmonitor/internal/auth"
	"gopkg.in/guregu/null.v3"
)

type controller struct {
	app       app.App
	templates *Templates
}

type Controller interface {
	Update(http.ResponseWriter, *http.Request)
}

func Init(app app.App) Controller {
	c := &controller{
		app:       app,
		templates: InitTemplates(),
	}

	auth, ok := app.Controller("auth").(auth.Controller)
	if !ok {
		c.app.Fatal("missing controller: auth")
	}

	_ = auth

	r := c.app.Router().PathPrefix("/computers").Subrouter()
	r.Handle("", http.HandlerFunc(c.Index)).Methods("GET").Name("computer_index")
	r.HandleFunc("/update", c.Update).Methods("POST").Name("update")
	r.HandleFunc("/{name}", c.Get).Methods("GET").Name("get")

	return c
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request) {

}

func (c *controller) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	compStore := NewComputerStore(c.app.DB())
	netStore := NewNetworkAdapterStore(c.app.DB())
	userStore := NewUserStore(c.app.DB())

	data, err := compStore.Get(r.Context(), vars["name"])
	if err != nil {
		c.app.Error(err)
	}

	if data != nil {
		netData, err := netStore.GetAllByComputerID(r.Context(), int(data.ID.Int64))
		if err != nil {
			c.app.Error(err)
		}

		data.NetworkAdapter = netData

		userData, err := userStore.GetAllUsersByComputerName(r.Context(), data.Name.String)
		if err != nil {
			c.app.Error(err)
		}

		data.Users = userData

		c.app.JsonResponse(w, data, 200)
		return
	}

}

func (c *controller) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 10))
	defer cancel()

	compStore := NewComputerStore(c.app.DB())
	netStore := NewNetworkAdapterStore(c.app.DB())
	userStore := NewUserStore(c.app.DB())

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.app.Error(err)
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
		c.app.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var compID null.Int

	comp, err := compStore.Get(ctx, record.Name.String)
	if err != nil {
		c.app.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comp != nil {

		// Computer record exists update the updated date field
		compID = null.IntFrom(comp.ID.Int64)
		err := compStore.UpdateLastActivityAt(ctx, comp)
		if err != nil {
			c.app.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {

		// Create new computer record
		compID, err = compStore.Create(ctx, &ComputerModel{
			Name: record.Name,
		})

		if err != nil {
			c.app.Error(err)
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
		c.app.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	networkAdapters, err := netStore.GetAllByComputerID(ctx, int(compID.Int64))
	if err != nil {
		c.app.Error(err)
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
						c.app.Error(err)
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
				c.app.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	}

	w.WriteHeader(http.StatusOK)
}

func (c *controller) Index(w http.ResponseWriter, r *http.Request) {
	compStore := NewComputerStore(c.app.DB())
	netStore := NewNetworkAdapterStore(c.app.DB())
	userStore := NewUserStore(c.app.DB())

	_ = netStore
	_ = userStore

	computers, err := compStore.GetAll(r.Context(), 0, 10)
	if err != nil {
		c.app.Error(err)
	}

	var data []ComputerModel
	for _, cp := range computers {
		n, err := netStore.GetAllByComputerID(r.Context(), int(cp.ID.Int64))
		if err != nil {
			c.app.Error(err)
		}

		u, err := userStore.GetAllUsersByComputerName(r.Context(), cp.Name.String)
		if err != nil {
			c.app.Error(err)
		}

		cp.NetworkAdapter = n
		cp.Users = u

		data = append(data, cp)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		c.app.Error(err)
	}

	fmt.Print(string(jsonData))

	c.templates.List(w, nil)

}

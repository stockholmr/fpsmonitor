package computer

import (
	"context"
	"encoding/json"
	"fpsmonitor/internal/logging"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

func Index(db *sqlx.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 10))
		defer cancel()

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logging.Error(err)
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
			logging.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		computerRepo := NewComputerRepository(db)
		userRepo := NewUserRepository(db)
		networkAdapterRepo := NewNetworkAdapterRepository(db)

		var compID int64

		comp, err := computerRepo.Select(ctx, record.Name.String)
		if err != nil {
			logging.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if comp != nil {

			// Computer record exists update the updated date field
			compID = comp.ID.Int64
			err := computerRepo.Update(ctx, comp)
			if err != nil {
				logging.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {

			// Create new computer record
			compID, err = computerRepo.Create(ctx, &Computer{
				Name: record.Name,
			})

			if err != nil {
				logging.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}

		user, err := userRepo.SelectWithUsernameAndComputerID(ctx, int(compID), record.Username.String)
		if err != nil {
			logging.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if user != nil {

			// User record exists
			err = userRepo.Update(ctx, user)
			if err != nil {
				logging.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {

			// Create new user record
			_, err = userRepo.Create(ctx, &User{
				Username:   record.Username,
				ComputerID: null.IntFrom(compID),
			})

			if err != nil {
				logging.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}

		networkAdapters, err := networkAdapterRepo.SelectWithComputerID(ctx, int(compID))
		if err != nil {
			logging.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(networkAdapters) > 0 {

			for _, na := range networkAdapters {
				for _, nar := range record.Adapters {
					if na.MacAddress == nar.MacAddress {
						na.IPAddress = nar.IPAddress
						na.Name = nar.Name
						err = networkAdapterRepo.Update(ctx, &na)
						if err != nil {
							logging.Error(err)
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
				}
			}

		} else {

			for _, na := range record.Adapters {
				na.ComputerID = null.IntFrom(compID)
				if _, err = networkAdapterRepo.Create(ctx, &na); err != nil {
					logging.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

		}

		w.WriteHeader(http.StatusOK)
	}
}

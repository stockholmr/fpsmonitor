package main

import (
	"context"
	"encoding/json"
	"fpsmonitor/internal/computer"
	"fpsmonitor/internal/logging"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v3"
)

func index(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 10))
	defer cancel()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var record struct {
		Name     null.String               `json:"name"`
		Username null.String               `json:"username"`
		Adapters []computer.NetworkAdapter `json:"adapters"`
	}

	err = json.Unmarshal(data, &record)
	if err != nil {
		logging.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	computerRepo := computer.NewComputerRepository(db)
	userRepo := computer.NewUserRepository(db)
	networkAdapterRepo := computer.NewNetworkAdapterRepository(db)

	existingRecord, err := computerRepo.Select(ctx, record.Name.String)
	if err != nil {
		logging.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if existingRecord != nil {

		err := computerRepo.Update(ctx, existingRecord)
		if err != nil {
			logging.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/*
			userRecord, err := userRepo.GetWithComputerID(ctx, record.Username.String, int(existingRecord.ID.Int64))
			if err != nil {
				logging.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if userRecord != nil {
				err = userRepo.Update(ctx, userRecord)
				if err != nil {
					logging.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}*/

		/*
			networkAdapters, err := networkAdapterRepo.Get(ctx, int(existingRecord.ID.Int64))
			if err != nil {
				logging.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if networkAdapters != nil {
				for _, na := range *networkAdapters {
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
			}
		*/

	} else {

		pc := computer.Computer{
			Name: record.Name,
		}
		pcID, err := computerRepo.Create(ctx, &pc)

		if err != nil {
			logging.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := computer.User{
			Username:   record.Username,
			ComputerID: null.IntFrom(pcID),
		}

		if _, err = userRepo.Create(ctx, &user); err != nil {
			logging.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, na := range record.Adapters {
			na.ComputerID = null.IntFrom(pcID)
			if _, err = networkAdapterRepo.Create(ctx, &na); err != nil {
				logging.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)

} // end index()

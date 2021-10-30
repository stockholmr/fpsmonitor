package main

import (
	"context"
	"encoding/json"
	"fpsmonitor/internal/computer"
	"fpsmonitor/internal/logging"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/guregu/null.v3"
)

var (
	listen = ""
	port   = "8080"
	router *mux.Router
	newDB  = false
	dbFile = "fpsmonitor.sqlite"
	db     *sqlx.DB
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Tracef("%s|%s|%s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

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

	existingRecord, err := computerRepo.Get(ctx, record.Name.String)
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
		}

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

func main() {

	// = Init Logger =========================================================================

	if err := logging.Init("fpsmonitor.log", logging.TRACE, logging.DEBUG); err != nil {
		logging.Debug(err)
	}

	// = Init Datebase Connection =========================================================================

	dbDir := path.Dir(dbFile)
	err := os.Mkdir(dbDir, 0776)
	if err != nil {
		if !os.IsExist(err) {
			logging.Fatalf("database failed: %s", err)
		}
	}

	if _, err := os.Stat(dbFile); err != nil {
		if !os.IsExist(err) {
			newDB = true
		}
	}

	db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		logging.Fatalf("database failed: %s", err)
	}

	dbCtx := context.Background()

	if newDB {
		if err = computer.NewComputerRepository(db).Install(dbCtx); err != nil {
			logging.Fatalf("database failed: %s", err)
		}

		if err = computer.NewNetworkAdapterRepository(db).Install(dbCtx); err != nil {
			logging.Fatalf("database failed: %s", err)
		}

		if err = computer.NewUserRepository(db).Install(dbCtx); err != nil {
			logging.Fatalf("database failed: %s", err)
		}
	}

	defer db.Close()

	// = Init Mux Router =========================================================================

	router = mux.NewRouter().StrictSlash(true)

	router.Handle("/", alice.New(LoggingMiddleware).ThenFunc(index)).Methods("POST")

	// = Init HTTP Server =========================================================================

	server := http.Server{
		Addr:           listen + ":" + port,
		Handler:        router,
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logging.Debugf("server started listening on address: %s port: %s", listen, port)
		err := server.ListenAndServe()
		if err != nil {
			logging.Fatalf("server failed: %s", err)
		}
	}()

	wait := time.Second * 15
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	logging.Info("server shutdown")
	os.Exit(0)
}

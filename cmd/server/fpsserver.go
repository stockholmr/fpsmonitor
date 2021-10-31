package main

import (
	"context"
	"fpsmonitor/internal/computer"
	"fpsmonitor/internal/logging"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	_ "github.com/mattn/go-sqlite3"
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

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/stockholmr/auth"
	"github.com/stockholmr/fpsmonitor/internal/assets"
	"github.com/stockholmr/fpsmonitor/internal/computer"
	"github.com/stockholmr/fpsmonitor/internal/logging"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	ini "gopkg.in/ini.v1"
)

/*
func (c *authController) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.log.Trace("%s|%s|%s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
*/

var (
	Server = struct {
		ListenAddress string `ini:"Listen"`
		Port          string `ini:"Port"`
		SessionKey    string `ini:"SessionKey"`
	}{
		ListenAddress: "",
		Port:          "8080",
	}

	Database = struct {
		File    string `ini:"File"`
		Install bool   `ini:"-"`
	}{
		File:    "fpsmonitor.sqlite",
		Install: false,
	}

	Logging = struct {
		File string `ini:"File"`
	}{
		File: "fpsmonitor.log",
	}

	configFile = "fpsmonitor.ini"
	router     *mux.Router
	db         *sqlx.DB
	cfg        *ini.File
)

func main() {

	_, err := os.Stat(configFile)
	if err != nil {
		if !os.IsExist(err) {
			cfg = ini.Empty()
			secServer, _ := cfg.NewSection("Server")
			secServer.ReflectFrom(&Server)
			secDatabase, _ := cfg.NewSection("Database")
			secDatabase.ReflectFrom(&Database)
			secLogging, _ := cfg.NewSection("Logging")
			secLogging.ReflectFrom(&Logging)
			cfg.SaveTo(configFile)
		}
	}

	cfg, err = ini.Load(configFile)
	if err != nil {
		logging.Fatal("failed to load config")
	}

	cfg.Section("Server").MapTo(&Server)
	cfg.Section("Database").MapTo(&Database)
	cfg.Section("Logging").MapTo(&Logging)

	// = Init Logger =========================================================================

	logger, err := logging.NewLogger(Logging.File, logging.TRACE, logging.DEBUG)
	if err != nil {
		logging.Debug(err)
	}

	// = Init Datebase Connection =========================================================================

	dbDir := path.Dir(Database.File)
	err = os.Mkdir(dbDir, 0776)
	if err != nil {
		if !os.IsExist(err) {
			logging.Fatalf("database failed: %s", err)
		}
	}

	if _, err := os.Stat(Database.File); err != nil {
		if !os.IsExist(err) {
			Database.Install = true
		}
	}

	db, err = sqlx.Open("sqlite3", Database.File)
	if err != nil {
		logging.Fatalf("database failed: %s", err)
	}

	dbCtx := context.Background()

	if err = auth.NewUserRepository(db).Install(); err != nil {
		logging.Fatalf("database failed: %s", err)
	}

	if Database.Install {
		if err = auth.NewUserRepository(db).Install(); err != nil {
			logging.Fatalf("database failed: %s", err)
		}

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

	// = Init Session Store ======================================================================

	sessionStore := sessions.NewCookieStore([]byte(Server.SessionKey))

	// = Init Mux Router =========================================================================

	router = mux.NewRouter().StrictSlash(true)

	router.Handle("/bootstrap", assets.Bootstrap()).Methods("GET")
	router.Handle("/jquery", assets.Jquery()).Methods("GET")
	router.Handle("/axios", assets.Axios()).Methods("GET")

	_ = auth.NewAuthController(db, logger, router, sessionStore, "list")
	_ = computer.NewComputerController(db, logger, router)

	//router.Handle("/", alice.New(LoggingMiddleware).ThenFunc(computer.Index(db))).Methods("POST")

	//router.Handle("/computers", alice.New(LoggingMiddleware).ThenFunc(computer.List(db))).Methods("GET", "POST")
	//router.Handle("/computers/stylesheet", computer.Stylesheet()).Methods("GET")

	// = Init HTTP Server =========================================================================

	server := http.Server{
		Addr:           Server.ListenAddress + ":" + Server.Port,
		Handler:        router,
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logging.Debugf("server started listening on address: %s port: %s", Server.ListenAddress, Server.Port)
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

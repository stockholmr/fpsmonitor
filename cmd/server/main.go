package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/stockholmr/fpsmonitor/internal/admin"
	"github.com/stockholmr/fpsmonitor/internal/auth"
	"github.com/stockholmr/fpsmonitor/internal/computer"
)

func main() {

	conf := InitConfig("fpsmonitor.ini")

	logg := InitLog(conf.Logging.File)

	db := InitDB(conf.Database.File, logg)
	defer db.Close()

	router := mux.NewRouter().StrictSlash(true)

	var sessionKeys *auth.SessionKeys
	sessionKeys, err := auth.SessionKeysFromBase64String(conf.Session.AuthenticationKey, conf.Session.EncryptionKey)
	if err != nil {
		sessionKeys.AuthenticationKey = securecookie.GenerateRandomKey(64)
		sessionKeys.EncryptionKey = securecookie.GenerateRandomKey(32)
	}

	authController := auth.InitWithLogger(router, db, sessionKeys, logg)
	computer.InitWithLogger(router, db, logg)
	admin.InitWithLogger(router, db, authController, logg)

	server := http.Server{
		Addr:           conf.Server.ListenAddress + ":" + conf.Server.Port,
		Handler:        router,
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logg.Printf("[HTTP SERVER] Address:%s Port:%s\n", conf.Server.ListenAddress, conf.Server.Port)
		err := server.ListenAndServe()
		if err != nil {
			logg.Println("[HTTP SERVER]", err)
		}
	}()

	wait := time.Second * 15
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	os.Exit(0)
}

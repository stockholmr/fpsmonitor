package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/stockholmr/fpsmonitor/internal/admin"
	"github.com/stockholmr/fpsmonitor/internal/app"
	"github.com/stockholmr/fpsmonitor/internal/auth"
	"github.com/stockholmr/fpsmonitor/internal/computer"
)

func main() {

	serverApp := app.New()

	serverApp.InitConfig("fpsmonitor.ini")
	serverApp.InitFileLog(serverApp.Config().Logging.File)
	serverApp.InitDB(serverApp.Config().Database.File)
	serverApp.InitSessionKeysFromBase64(
		serverApp.Config().Session.AuthenticationKey,
		serverApp.Config().Session.EncryptionKey,
	)
	serverApp.InitRouter()

	defer serverApp.DB().Close()

	authController := auth.InitWithLogger(router, db, sessionKeys, logg)
	computer.InitWithLogger(router, db, logg)
	admin.InitWithLogger(router, db, authController, logg)

	server := http.Server{
		Addr:           serverApp.Config().Server.ListenAddress + ":" + serverApp.Config().Server.Port,
		Handler:        serverApp.Router(),
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		serverApp.Info(
			"[HTTP SERVER]",
			"Address:", serverApp.Config().Server.ListenAddress,
			"Port:", serverApp.Config().Server.Port,
		)
		err := server.ListenAndServe()
		if err != nil {
			serverApp.Error("[HTTP SERVER]", err)
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

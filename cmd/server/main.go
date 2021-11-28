package main

import (
	"os"
	"os/signal"

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
	serverApp.InitCsrfMiddleware()
	serverApp.InitRouter()

	defer serverApp.DB().Close()

	serverApp.RegisterController("auth", auth.Init(serverApp))
	serverApp.RegisterController("admin", admin.Init(serverApp))
	serverApp.RegisterController("computer", computer.Init(serverApp))

	serverApp.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	serverApp.Stop()
	os.Exit(0)
}

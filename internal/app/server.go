package app

import (
	"context"
	"net/http"
	"time"
)

func (a *App) Run() {

	a.server = &http.Server{
		Addr:           a.Config().Server.ListenAddress + ":" + a.Config().Server.Port,
		Handler:        a.Router(),
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		a.Info(
			"[HTTP SERVER]",
			"Address:", a.Config().Server.ListenAddress,
			"Port:", a.Config().Server.Port,
		)
		err := a.server.ListenAndServe()
		if err != nil {
			a.Error("[HTTP SERVER]", err)
		}
	}()
}

func (a *App) Stop() {
	wait := time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	a.server.Shutdown(ctx)
}

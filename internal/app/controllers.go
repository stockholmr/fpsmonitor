package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *App) RegisterController(name string, c interface{}) {
	for cName, _ := range a.controllers {
		if cName == name {
			a.Fatal("controller already exists with name:", name)
		}
	}
	a.controllers[name] = c
}

func (a *App) Controller(name string) interface{} {
	for cName, c := range a.controllers {
		if cName == name {
			return c
		}
	}
	a.Fatal("invalid controller name")
	return nil
}

func (a *App) JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonStr, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		a.Error(err)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonStr)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(jsonStr))
}

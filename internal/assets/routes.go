package assets

import (
	"fmt"
	"net/http"
)

func Axios() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(axios())))
		w.Write(axios())
	}
}

func Bootstrap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bootstrap())))
		w.Write(bootstrap())
	}
}

func Jquery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jquery())))
		w.Write(jquery())
	}
}

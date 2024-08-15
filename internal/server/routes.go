package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func addRoutes(mux *httprouter.Router) {
	//Authentication Endpoints
	mux.GET("/create-account", func(f http.ResponseWriter, r *http.Request, p httprouter.Params) {
		json.NewEncoder(f).Encode("Hello World")
	})
}

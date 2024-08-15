package server

import (
	"github.com/julienschmidt/httprouter"
)

func addRoutes(mux *httprouter.Router) {
  //Authentication Endpoints
	mux.POST("/create-account", auth.CreateAccountHandler)
}

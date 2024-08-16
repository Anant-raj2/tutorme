package server

import (
	"github.com/julienschmidt/httprouter"
)

func (srv *HTTP) addRoutes(mux *httprouter.Router) {
	//Authentication Endpoints
	mux.POST("/create-account", ErrorWrapper(srv.handleCreateAccount))
}

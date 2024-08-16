package server

import (
	"github.com/julienschmidt/httprouter"
)

func (srv *HTTP) addRoutes(mux *httprouter.Router) {
	//Authentication Endpoints
	mux.GET("/signup", srv.userHandler.renderSignup)
	mux.POST("/create-account", ErrorWrapper(srv.userHandler.CreateTutor))
}

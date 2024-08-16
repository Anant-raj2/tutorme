package server

import (
	"github.com/Anant-raj2/tutorme/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

func (srv *HTTP) addRoutes(mux *httprouter.Router) {
	//Authentication Endpoints
	mux.GET("/signup", srv.authStore.RenderSignup)
	mux.POST("/create-account", logger.LogHttp(ErrorWrapper(srv.authStore.CreateTutor)))
}

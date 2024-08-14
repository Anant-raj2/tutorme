package server

import (
	"github.com/Anant-raj2/tutorme/internal/handler"
	"github.com/julienschmidt/httprouter"
)

func addRoutes(mux *httprouter.Router) {
	mux.GET("/data", handler.HelloWorld)
}

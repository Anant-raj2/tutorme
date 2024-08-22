package server

import (
	"context"
	"net/http"

	"github.com/Anant-raj2/tutorme/pkg/logger"
	"github.com/Anant-raj2/tutorme/web/templa/home"
	"github.com/julienschmidt/httprouter"
)

func (srv *HTTP) addRoutes(mux *httprouter.Router) {
	// Home Endpoints
	mux.GET("/", logger.LogHttp(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		component := home.HomeLayout()
		component.Render(context.Background(), w)
	}))

	//Tutor Endpoints
	mux.GET("/register", logger.LogHttp(srv.authHandler.RenderRegister))
	mux.POST("/create/account", logger.LogHttp(ErrorWrapper(srv.authHandler.Register)))

	//Tutor Endpoints
	mux.GET("/create-tutor", logger.LogHttp(srv.tutorHandler.RenderTutor))
	mux.POST("/create/tutor", logger.LogHttp(ErrorWrapper(srv.tutorHandler.CreateTutor)))
}

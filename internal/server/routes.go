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
	//Authentication Endpoints
	mux.GET("/signup", logger.LogHttp(srv.authStore.RenderSignup))
	mux.POST("/create-account", logger.LogHttp(ErrorWrapper(srv.authStore.CreateTutor)))
}

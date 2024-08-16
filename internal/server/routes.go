package server

import (
	// "context"
	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srv *HTTP) addRoutes(mux *httprouter.Router) {
	//Authentication Endpoints
	mux.POST("/create-account", ErrorWrapper(srv.handleCreateAccount))
}

func (srv *HTTP) handleCreateAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	// ctx := context.Background()
	fmt.Println(r.Body)
	// srv.queries.CreateTutor(ctx, arg db.TutorParams)
	return nil
}

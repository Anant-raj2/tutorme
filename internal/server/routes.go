package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/julienschmidt/httprouter"
)

func (srv *HTTP) addRoutes(mux *httprouter.Router) {
	//Authentication Endpoints
	mux.POST("/create-account", ErrorWrapper(srv.handleCreateAccount))
	mux.GET("/test-account", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		json.NewEncoder(w).Encode("Hello")
	})
}

type Sentence struct {
	Sentence string `json:"sentence"`
}

func (srv *HTTP) handleCreateAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	ctx := context.Background()
	user := db.TutorParams{}

	req, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(req, &user)

	tutor, err := srv.queries.CreateTutor(ctx, user)
	if err != nil {
		return err
	}
  fmt.Printf("%s", tutor)
	json.NewEncoder(w).Encode("hello")
	return nil
}

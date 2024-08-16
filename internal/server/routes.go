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

type TutorParams struct {
	Name       string `json:"name"`
	GradeLevel int32  `json:"grade_level"`
	Role       string `json:"role"`
	Gender     string `json:"gender"`
	Subject    string `json:"subject"`
}

func (params *TutorParams) OTD() *db.CreateTutorParams {
	return &db.CreateTutorParams{
		Name:       params.Name,
		GradeLevel: params.GradeLevel,
		Role:       params.Role,
		Gender:     params.Gender,
		Subject:    params.Subject,
	}
}

func (srv *HTTP) handleCreateAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	ctx := context.Background()
	var user TutorParams

	req, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(req, &user)

	tutor, err := srv.queries.CreateTutor(ctx, *user.OTD())
	if err != nil {
		return err
	}
	fmt.Println(tutor)
	json.NewEncoder(w).Encode(tutor)
	return nil
}

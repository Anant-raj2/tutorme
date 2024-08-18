package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/Anant-raj2/tutorme/web/templa/auth"
	"github.com/julienschmidt/httprouter"
)

type AuthStore struct {
	queries *db.Queries
}

func New(queries *db.Queries) *AuthStore {
	return &AuthStore{
		queries: queries,
	}
}

func (handler *AuthStore) RenderSignup(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	component := auth.SignupForm("Anant-raj2")
	component.Render(context.Background(), w)
}

func (handler *AuthStore) CreateTutor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	ctx := context.Background()
	r.ParseForm()
	grade_level, err := strconv.Atoi(r.FormValue("grade_level"))

	var tutorConfig db.CreateTutorParams = db.CreateTutorParams{
		Email:      r.FormValue("email"),
		Name:       r.FormValue("name"),
		Gender:     r.FormValue("gender"),
		GradeLevel: int32(grade_level),
		Role:       r.FormValue("role"),
		Subject:    r.FormValue("subject"),
	}

	tutor, err := handler.queries.CreateTutor(ctx, tutorConfig)
	if err != nil {
		return err
	}

	fmt.Println(tutor)
	component := auth.SignupForm("Anant-raj2")
	component.Render(context.Background(), w)
	return nil
}

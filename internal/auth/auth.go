package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/Anant-raj2/tutorme/pkg/util"
	"github.com/Anant-raj2/tutorme/web/templa"
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
	component := templa.Hello("Anant-raj2")
	component.Render(context.Background(), w)
}

func (handler *AuthStore) CreateTutor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	ctx := context.Background()
	var user util.TutorParams

	req, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(req, &user)

	start := time.Now()

	tutor, err := handler.queries.CreateTutor(ctx, *user.OTD())
	if err != nil {
		return err
	}

	fmt.Println(time.Since(start).Milliseconds())

	fmt.Println(user)

	json.NewEncoder(w).Encode(tutor)
	return nil
}

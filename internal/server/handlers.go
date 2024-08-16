package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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

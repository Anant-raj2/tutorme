package auth

import (
	"context"
	"net/http"

	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/Anant-raj2/tutorme/web/templa/auth"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	*db.Queries
}

func New(db *db.Queries) *Handler {
	return &Handler{
		db,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	return nil
}

func (h *Handler) RenderRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	component := auth.Register()
	component.Render(context.Background(), w)
}

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anant-raj2/tutorme/internal/db"
)

type TutorParams struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	GradeLevel int32  `json:"grade_level"`
	Role       string `json:"role"`
	Gender     string `json:"gender"`
	Subject    string `json:"subject"`
}

func (params *TutorParams) OTD() *db.CreateTutorParams {
	return &db.CreateTutorParams{
		Name:       params.Name,
		Email:      params.Email,
		GradeLevel: params.GradeLevel,
		Role:       params.Role,
		Gender:     params.Gender,
		Subject:    params.Subject,
	}
}

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

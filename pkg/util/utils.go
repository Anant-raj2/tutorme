package util

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
	Gender     string `json:"gender"`
	Subject    string `json:"subject"`
}

func (params *TutorParams) OTD() *db.CreateTutorParams {
	return &db.CreateTutorParams{
		Name:       params.Name,
		Email:      params.Email,
		GradeLevel: params.GradeLevel,
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

func (list *LinkedList) findNodeAt(index int) *Node {
 var count int = 0
 var current *Node = list.head

 for current != nil {
  count++
  current = current.next
 }

 if index <= 0 || index > count {
  return nil
 }

 current = list.head
 for count = 1; count < index; count++ {
  current = current.next
 }
 return current
}


func (list *LinkedList) print() {
 var current *Node = list.head
 for current != nil {
  fmt.Printf("%d -> ", current.data)
  current = current.next
 }
 fmt.Println()
}

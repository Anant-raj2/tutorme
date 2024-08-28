package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Tutor struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Subject      string    `json:"subject"`
	Experience   int       `json:"experience"`
	Rate         float64   `json:"rate"`
	Availability []string  `json:"availability"`
	CreatedAt    time.Time `json:"created_at"`
}

var pool *pgxpool.Pool

func main() {
	var err error
	pool, err = pgxpool.Connect(context.Background(), "postgres://username:password@localhost/tutordb")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/tutors", createTutor).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createTutor(w http.ResponseWriter, r *http.Request) {
	var tutor Tutor
	err := json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	sqlStatement := `
	INSERT INTO tutors (name, subject, experience, rate, availability, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at`

	err = tx.QueryRow(ctx, sqlStatement, tutor.Name, tutor.Subject, tutor.Experience, tutor.Rate, tutor.Availability, time.Now()).Scan(&tutor.ID, &tutor.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutor)
}

// handle a CRUD request
func (rtr *CRUDrouter[T]) handle(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.EscapedPath())
	// We need ID for anything other than a GET to muxbase (AKA list request)
	if r.URL.EscapedPath() != rtr.muxBase && err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Instantiate zero-value for our model
	var model T
	switch r.Method {
	case http.MethodGet:
		// List items
		if id == 0 {
			items := model.GetList()
			json.NewEncoder(w).Encode(&items)
		} else {
			// Get single item
			item := model.GetSingleItem(id)
			json.NewEncoder(w).Encode(&item)
		}
	case http.MethodPost:
		// Create a new record.
		json.NewDecoder(r.Body).Decode(&model)
		model.Create()
	case http.MethodPut:
		// Update an existing record.
		json.NewDecoder(r.Body).Decode(&model)
		model.UpdateItem(id)
	case http.MethodDelete:
		// Remove the record.
		model.Delete(id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

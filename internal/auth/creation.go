package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type Tutor struct {
	ID           int      `json:"id" db:"id"`
	Name         string   `json:"name" db:"name"`
	Subject      string   `json:"subject" db:"subject"`
	Experience   int      `json:"experience" db:"experience"`
	Rate         float64  `json:"rate" db:"rate"`
	Availability []string `json:"availability" db:"availability"`
}

var db *sqlx.DB

func main() {
	var err error
	db, err = sqlx.Connect("postgres", "postgres://username:password@localhost/tutordb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	tx, err := db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement := `
	INSERT INTO tutors (name, subject, experience, rate, availability)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	err = tx.QueryRowx(sqlStatement, tutor.Name, tutor.Subject, tutor.Experience, tutor.Rate, tutor.Availability).Scan(&tutor.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutor)
}

package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Tutor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Subject   string `json:"subject"`
	Experience int   `json:"experience"`
	Rate      float64 `json:"rate"`
}

func setupCreation() {
	var err error
	db, err = sql.Open("postgres", "postgres://username:password@localhost/tutordb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}

func createTutor(w http.ResponseWriter, r *http.Request) {
	var tutor Tutor
	err := json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `
	INSERT INTO tutors (name, subject, experience, rate)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	err = db.QueryRow(sqlStatement, tutor.Name, tutor.Subject, tutor.Experience, tutor.Rate).Scan(&tutor.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tutor)
}

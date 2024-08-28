package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Student struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Age            int       `json:"age"`
	Grade          int       `json:"grade"`
	SubjectsNeeded []string  `json:"subjects_needed"`
	CreatedAt      time.Time `json:"created_at"`
}

var db *sql.DB

func studentdb() {
	var err error
	db, err = sql.Open("postgres", "postgres://username:password@localhost/studentdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/students", createStudent).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	sqlStatement := `
	INSERT INTO students (name, email, age, grade, subjects_needed, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at`

	err = tx.QueryRow(sqlStatement, student.Name, student.Email, student.Age, student.Grade, pq.Array(student.SubjectsNeeded), time.Now()).Scan(&student.ID, &student.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

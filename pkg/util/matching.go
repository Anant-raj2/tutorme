package util

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Tutor struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Subjects   []string `json:"subjects"`
	Experience int      `json:"experience"`
	Rate       float64  `json:"rate"`
}

type Student struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	SubjectsNeeded []string `json:"subjects_needed"`
	Grade          int      `json:"grade"`
}

type Match struct {
	TutorID   int     `json:"tutor_id"`
	StudentID int     `json:"student_id"`
	Score     float64 `json:"score"`
}

var pool *pgxpool.Pool

func main() {
	var err error
	pool, err = pgxpool.Connect(context.Background(), "postgres://username:password@localhost/tutormatchdb")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/match", matchTutorsWithStudents).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func matchTutorsWithStudents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tutors, err := fetchTutors(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch tutors", http.StatusInternalServerError)
		return
	}

	students, err := fetchStudents(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}

	matches := performMatching(tutors, students)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func fetchTutors(ctx context.Context) ([]Tutor, error) {
	rows, err := pool.Query(ctx, "SELECT id, name, subjects, experience, rate FROM tutors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tutors []Tutor
	for rows.Next() {
		var t Tutor
		err := rows.Scan(&t.ID, &t.Name, &t.Subjects, &t.Experience, &t.Rate)
		if err != nil {
			return nil, err
		}
		tutors = append(tutors, t)
	}

	return tutors, nil
}

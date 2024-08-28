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

func fetchStudents(ctx context.Context) ([]Student, error) {
	rows, err := pool.Query(ctx, "SELECT id, name, subjects_needed, grade FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.ID, &s.Name, &s.SubjectsNeeded, &s.Grade)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}

func performMatching(tutors []Tutor, students []Student) []Match {
	var matches []Match

	for _, student := range students {
		for _, tutor := range tutors {
			score := calculateMatchScore(tutor, student)
			if score > 0 {
				matches = append(matches, Match{
					TutorID:   tutor.ID,
					StudentID: student.ID,
					Score:     score,
				})
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches[:min(len(matches), 10)] // Return top 10 matches
}

func calculateMatchScore(tutor Tutor, student Student) float64 {
	var score float64

	for _, subject := range student.SubjectsNeeded {
		for _, tutorSubject := range tutor.Subjects {
			if subject == tutorSubject {
				score += 1.0
				break
			}
		}
	}

	if score == 0 {
		return 0
	}

	experienceBonus := float64(tutor.Experience) * 0.1
	score += experienceBonus

	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

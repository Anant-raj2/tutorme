package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Tutor struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Subjects     []string  `json:"subjects"`
	Experience   int       `json:"experience"`
	Rate         float64   `json:"rate"`
	Availability []string  `json:"availability"`
	LastActive   time.Time `json:"last_active"`
}

type Student struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	SubjectsNeeded []string  `json:"subjects_needed"`
	Grade          int       `json:"grade"`
	PreferredTimes []string  `json:"preferred_times"`
	LastActive     time.Time `json:"last_active"`
}

type Match struct {
	TutorID          int     `json:"tutor_id"`
	StudentID        int     `json:"student_id"`
	Score            float64 `json:"score"`
	MatchedSubjects  []string `json:"matched_subjects"`
	AvailabilityMatch float64 `json:"availability_match"`
}

var pool *pgxpool.Pool

func mathching() {
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
	rows, err := pool.Query(ctx, "SELECT id, name, subjects, experience, rate, availability, last_active FROM tutors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tutors []Tutor
	for rows.Next() {
		var t Tutor
		err := rows.Scan(&t.ID, &t.Name, &t.Subjects, &t.Experience, &t.Rate, &t.Availability, &t.LastActive)
		if err != nil {
			return nil, err
		}
		tutors = append(tutors, t)
	}

	return tutors, nil
}

func fetchStudents(ctx context.Context) ([]Student, error) {
	rows, err := pool.Query(ctx, "SELECT id, name, subjects_needed, grade, preferred_times, last_active FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.ID, &s.Name, &s.SubjectsNeeded, &s.Grade, &s.PreferredTimes, &s.LastActive)
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
			score, matchedSubjects, availabilityMatch := calculateMatchScore(tutor, student)
			if score > 0 {
				matches = append(matches, Match{
					TutorID:          tutor.ID,
					StudentID:        student.ID,
					Score:            score,
					MatchedSubjects:  matchedSubjects,
					AvailabilityMatch: availabilityMatch,
				})
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches[:min(len(matches), 10)] // Return top 10 matches
}

func calculateMatchScore(tutor Tutor, student Student) (float64, []string, float64) {
	var score float64
	var matchedSubjects []string
	var availabilityMatch float64

	// Subject matching with fuzzy search
	for _, studentSubject := range student.SubjectsNeeded {
		for _, tutorSubject := range tutor.Subjects {
			if fuzzy.Match(strings.ToLower(studentSubject), strings.ToLower(tutorSubject)) {
				score += 1.0
				matchedSubjects = append(matchedSubjects, studentSubject)
				break
			}
		}
	}

	if len(matchedSubjects) == 0 {
		return 0, nil, 0
	}

	// Experience bonus
	experienceBonus := float64(tutor.Experience) * 0.1
	score += experienceBonus

	// Availability matching
	availabilityMatch = calculateAvailabilityMatch(tutor.Availability, student.PreferredTimes)
	score += availabilityMatch

	// Recency bonus
	recencyBonus := calculateRecencyBonus(tutor.LastActive, student.LastActive)
	score += recencyBonus

	return score, matchedSubjects, availabilityMatch
}

func calculateAvailabilityMatch(tutorAvailability, studentPreferredTimes []string) float64 {
	var matchCount int
	for _, tutorTime := range tutorAvailability {
		for _, studentTime := range studentPreferredTimes {
			if tutorTime == studentTime {
				matchCount++
				break
			}
		}
	}
	return float64(matchCount) / float64(len(studentPreferredTimes))
}

func calculateRecencyBonus(tutorLastActive, studentLastActive time.Time) float64 {
	tutorDays := time.Since(tutorLastActive).Hours() / 24
	studentDays := time.Since(studentLastActive).Hours() / 24

	tutorBonus := 1.0 / (1.0 + tutorDays/30) // Decay over a month
	studentBonus := 1.0 / (1.0 + studentDays/30)

	return (tutorBonus + studentBonus) / 2
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

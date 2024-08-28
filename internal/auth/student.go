package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Student struct {
	ID             int       `json:"id"`
	Name           string    `json:"name" validate:"required,min=2,max=100"`
	Email          string    `json:"email" validate:"required,email"`
	Age            int       `json:"age" validate:"required,gte=5,lte=100"`
	Grade          int       `json:"grade" validate:"required,gte=1,lte=12"`
	SubjectsNeeded []string  `json:"subjects_needed" validate:"required,min=1,dive,required"`
	CreatedAt      time.Time `json:"created_at"`
}

var pool *pgxpool.Pool
var validate *validator.Validate

func studentdb() {
	var err error
	pool, err = pgxpool.Connect(context.Background(), "postgres://username:password@localhost/studentdb")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	validate = validator.New()
	err = validate.RegisterValidation("subjectformat", validateSubjectFormat)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/students", createStudent).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validate.Struct(student)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, e := range validationErrors {
			errorMessages[e.Field()] = e.Tag()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessages)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	sqlStatement := `
	INSERT INTO students (name, email, age, grade, subjects_needed, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at`

	err = tx.QueryRow(ctx, sqlStatement, student.Name, student.Email, student.Age, student.Grade, student.SubjectsNeeded, time.Now()).Scan(&student.ID, &student.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func validateSubjectFormat(fl validator.FieldLevel) bool {
	subject := fl.Field().String()
	match, _ := regexp.MatchString(`^[A-Z][a-z]+([ ][A-Z][a-z]+)*$`, subject)
	return match
}

func getStudentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var student Student
	err := pool.QueryRow(ctx, "SELECT id, name, email, age, grade, subjects_needed, created_at FROM students WHERE id = $1", id).
		Scan(&student.ID, &student.Name, &student.Email, &student.Age, &student.Grade, &student.SubjectsNeeded, &student.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

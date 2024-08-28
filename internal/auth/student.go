package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name           string    `json:"name"`
	Email          string    `json:"email" gorm:"uniqueIndex"`
	Age            int       `json:"age"`
	Grade          int       `json:"grade"`
	SubjectsNeeded []string  `json:"subjects_needed" gorm:"type:text[]"`
	CreatedAt      time.Time `json:"created_at"`
}

var db *gorm.DB

func studentdb() {
	var err error
	dsn := "host=localhost user=username password=password dbname=studentdb port=5432 sslmode=disable TimeZone=UTC"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Student{})
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	student.CreatedAt = time.Now()

	result := db.Create(&student)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func validateStudent(student *Student) error {
	if student.Name == "" {
		return errors.New("name is required")
	}
	if student.Email == "" {
		return errors.New("email is required")
	}
	if student.Age < 5 || student.Age > 100 {
		return errors.New("invalid age")
	}
	if student.Grade < 1 || student.Grade > 12 {
		return errors.New("invalid grade")
	}
	return nil
}

func getStudentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var student Student
	result := db.First(&student, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

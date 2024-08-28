package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/time/rate"
)

var (
	pool   *pgxpool.Pool
	rdb    *redis.Client
	limiter *rate.Limiter
)

type MatchRequest struct {
	StudentID int `json:"student_id"`
}

type MatchResponse struct {
	Matches []Match `json:"matches"`
}

func optim() {
	var err error

	// Initialize PostgreSQL connection pool
	pool, err = pgxpool.Connect(context.Background(), "postgres://username:password@localhost/tutormatchdb")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Initialize rate limiter: 100 requests per minute
	limiter = rate.NewLimiter(rate.Every(time.Minute/100), 100)

	router := mux.NewRouter()
	router.HandleFunc("/match", rateLimitMiddleware(cacheMiddleware(matchHandler))).Methods("POST")

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func cacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		r.Body.Close()

		cacheKey := fmt.Sprintf("match:%d", req.StudentID)

		// Try to get the result from cache
		cachedResult, err := rdb.Get(r.Context(), cacheKey).Result()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedResult))
			return
		}

		// If not in cache, call the handler
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		// Cache the result for 5 minutes
		if rec.Code == http.StatusOK {
			rdb.Set(r.Context(), cacheKey, rec.Body.String(), 5*time.Minute)
		}

		// Copy the response to the original ResponseWriter
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	}
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	var req MatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	student, err := fetchStudent(ctx, req.StudentID)
	if err != nil {
		http.Error(w, "Failed to fetch student", http.StatusInternalServerError)
		return
	}

	tutors, err := fetchTutors(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch tutors", http.StatusInternalServerError)
		return
	}

	matches := performMatching(tutors, *student)

	response := MatchResponse{Matches: matches}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

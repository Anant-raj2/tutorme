package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func main() {
	http.HandleFunc("/hello", HandleRoot)
	fmt.Println("Running server on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	var post Post = Post{
		ID:   1,
		Body: "Hello this is hello world",
	}
	fmt.Println("Sending Post...")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

package main

import (
	// "encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.GET("/send-data", SendData)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func SendData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
	// json.NewEncoder(w).Encode("Hello World")
}

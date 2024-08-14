package main

import (
	"log"

	"github.com/Anant-raj2/tutorme/internal/server"
)

func main() {
	srv, err := server.CreateService()
	if err != nil {
		log.Fatal("Could not create the server: ", err)
	}
	srv.Start()
}

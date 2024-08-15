package main

import (
	"github.com/Anant-raj2/tutorme/internal/server"
)

func main() {
	srv := server.NewHttpServer(server.DefaultConfig())
	srv.Start()
}

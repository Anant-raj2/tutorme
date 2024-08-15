package main

import (
	"context"

	"github.com/Anant-raj2/tutorme/internal/auth"
	"github.com/Anant-raj2/tutorme/internal/server"
)

func main() {
	var ctx context.Context = context.Background()

	srv := server.NewHttpServer(server.DefaultConfig(), authStore)

	srv.Start(ctx)
}

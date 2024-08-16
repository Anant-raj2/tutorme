package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Anant-raj2/tutorme/internal/auth"
	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/Anant-raj2/tutorme/internal/server"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load the environment file: %s", err)
	}

	var ctx context.Context = context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to database: %s", err)
	}

	defer conn.Close(ctx)

	queries := db.New(conn)

	var userHandler *auth.AuthStore = auth.New(queries)

	var handlerConfig server.Config = server.Config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
	var httpServer *server.HTTP = server.NewHttpServer(handlerConfig, userHandler)
	httpServer.Start(ctx)
}

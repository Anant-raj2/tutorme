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
	system, err := SetupEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load the environment")
		os.Exit(1)
	}

	var ctx context.Context = context.Background()

	conn, err := pgx.Connect(ctx, system.Db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to database: %s", err)
	}

	defer conn.Close(ctx)

	queries := db.New(conn)

	var userHandler *auth.AuthStore = auth.New(queries)

	var handlerConfig server.Config = server.Config{
		Host: system.Host,
		Port: system.Port,
	}
	var httpServer *server.HTTP = server.NewHttpServer(handlerConfig, userHandler)
	httpServer.Start(ctx)
}

type System struct {
	Host string
	Port string
	Db   string
}

func SetupEnv() (*System, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load the environment file: %s", err)
	}

	system := &System{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		Db:   os.Getenv("DB_URL"),
	}

	return system, err
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/Anant-raj2/tutorme/internal/server"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// connStr := os.Getenv("DB_URL")
	// conn, err := pgx.Connect(ctx, connStr)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not connect to database: %s", err)
	// }
	//
	// defer conn.Close(ctx)
	// queries := db.New(conn)
	//
	// tutorParams := db.TutorParams{
	// 	Name:       "Keith Decker",
	// 	Role:       "Admin",
	// 	Gender:     "male",
	// 	Subject:    "",
	// 	GradeLevel: 13,
	// }
	//
	// insertedAuthor, err := queries.CreateTutor(ctx, tutorParams)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not insert author: %s", err)
	// }
	// log.Println(insertedAuthor)
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load the environment file: %s", err)
	}

	var ctx context.Context = context.Background()

	var handlerConfig server.Config = server.Config{
		Host: os.Getenv("Host"),
		Port: os.Getenv("Port"),
	}

	var httpServer *server.HTTP = server.NewHttpServer(handlerConfig, nil)
	httpServer.Start(ctx)
}

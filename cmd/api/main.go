package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	var ctx context.Context = context.Background()
	connStr := os.Getenv("DB_URL")
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to database: %s", err)
	}

	defer conn.Close(ctx)
	db := db.New(conn)

	insertedAuthor, err := db.CreateTutor(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not insert author: %s", err)
	}
	log.Println(insertedAuthor)
}

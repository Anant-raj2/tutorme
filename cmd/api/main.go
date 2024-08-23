package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Anant-raj2/tutorme/internal/auth"
	"github.com/Anant-raj2/tutorme/internal/db"
	"github.com/Anant-raj2/tutorme/internal/server"
	"github.com/Anant-raj2/tutorme/internal/tutor"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	system, err := SetupEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load the environment")
		os.Exit(1)
	}
	// Create a logger with custom options
	logger := advancedlog.NewLogger(
		advancedlog.WithLevel(advancedlog.DEBUG),
		advancedlog.WithOutput(os.Stdout),
		advancedlog.WithFormatter(&advancedlog.JSONFormatter{}),
		advancedlog.WithHook(advancedlog.NewFileHook("app.log", 10, 5, 30)),
		advancedlog.WithFields(advancedlog.Field{Key: "app", Value: "myapp"}),
	)

	// Log some messages
	logger.Debug("This is a debug message")
	logger.Info("This is an info message", advancedlog.Field{Key: "user", Value: "john"})

	// Create a new logger with additional fields
	userLogger := logger.WithField("user", "alice")
	userLogger.Warn("This is a warning message")

	// Log an error
	err := someFunction()
	if err != nil {
		logger.Error("An error occurred", advancedlog.Field{Key: "error", Value: err.Error()})
	}

	// Fatal error
	logger.Fatal("This is a fatal error")

	var ctx context.Context = context.Background()

	conn, err := pgx.Connect(ctx, system.Db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to database: %s", err)
	}

	defer conn.Close(ctx)

	queries := db.New(conn)

	var handlerConfig server.Config = server.Config{
		Host: system.Host,
		Port: system.Port,
	}

	var userHandler *tutor.Handler = tutor.New(queries)
	var authHandler *auth.Handler = auth.New(queries)

	var httpServer *server.HTTP = server.NewHttpServer(handlerConfig, userHandler, authHandler)
	httpServer.Start(ctx)
}

type System struct {
	Host string
	Port string
	Db   string
}

func SetupEnv() (*System, error) {
	// Create a logger that writes to console
	consoleLogger := golog.NewLogger(golog.DEBUG, os.Stdout)

	// Log some messages
	consoleLogger.Debug("This is a debug message")
	consoleLogger.Info("This is an info message")
	consoleLogger.Warn("This is a warning message")
	consoleLogger.Error("This is an error message")

	// Create a logger that writes to a file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		consoleLogger.Fatal("Failed to open log file: %v", err)
	}
	defer file.Close()

	fileLogger := golog.NewLogger(golog.INFO, file)

	// Log some messages to the file
	fileLogger.Info("This message goes to the file")
	fileLogger.Error("An error occurred: %s", "Something went wrong")

	// Change the log level
	fileLogger.SetLevel(golog.WARN)
	fileLogger.Info("This message won't be logged because the level is set to WARN")
	fileLogger.Warn("This warning message will be logged")

	// Log to both console and file
	multiWriter := io.MultiWriter(os.Stdout, file)
	multiLogger := golog.NewLogger(golog.DEBUG, multiWriter)
	multiLogger.Info("This message goes to both console and file")
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

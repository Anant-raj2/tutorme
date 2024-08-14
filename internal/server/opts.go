package server

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type ServerConfig func(*Config)

func defaultConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load env")
		os.Exit(1)
	}
	return Config{
		Host:         os.Getenv("HOST"),
		Port:         os.Getenv("PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *pgxpool.Pool

func DBConnect() {
	err := LoadEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	dsn := "host=" + DB_HOST + " dbname=" + DB_NAME + " connect_timeout=5 user=" + DB_USER + " password=" + DB_PASSWORD + " sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	DB = pool
}

// Optional: Close the connection when your application shuts down
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func LoadEnv() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *pgxpool.Pool

func DBConnect() {
	dsn := "host=localhost dbname=mydb connect_timeout=5 user=postgres password=postgres sslmode=disable"

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

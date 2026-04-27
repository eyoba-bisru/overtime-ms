package config

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func DBConnect() {
	db, err := sql.Open("postgres", "host=localhost dbname=mydb connect_timeout=5 user=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Verify connection works
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db // Store the open connection
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

package main

import (
	"log"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	dotenvErr := config.LoadEnv()
	if dotenvErr != nil {
		panic(dotenvErr)
	}

	config.DBConnect()
	defer config.CloseDB()
	if err := config.Migrate(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.POST("/register", handlers.CreateUserHandler)
	r.POST("/login", handlers.LoginHandler)
	r.Run(":8080")
}

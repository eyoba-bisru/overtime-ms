package main

import (
	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	dotenvErr := config.LoadEnv()
	if dotenvErr != nil {
		panic(dotenvErr)
	}

	config.DBConnect()
	defer config.CloseDB()
	config.CreateTables()

	r := gin.Default()
	r.POST("/register", handlers.CreateUserHandler)
	r.POST("/login", handlers.LoginHandler)
	r.Run(":8080")
}

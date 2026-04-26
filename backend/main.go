package main

import (
	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config.DBConnect()
	defer config.CloseDB()
	config.CreateTables()

	r := gin.Default()
	r.POST("/register", handlers.CreateUserHandler)
	r.Run(":8080")
}

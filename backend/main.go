package main

import (
	"log"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/handlers"
	"github.com/eyoba-bisru/overtime-backend/internal/middleware"
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

	base := r.Group("/api/v1")

	// base.Use(middleware.CORSMiddleware())
	base.Use(middleware.LoggerMiddleware())

	auth := base.Group("/auth")
	auth.POST("/register", handlers.CreateUserHandler)
	auth.POST("/login", handlers.LoginHandler)

	overtime := base.Group("/overtime")
	overtime.Use(middleware.AuthMiddleware())
	overtime.POST("", handlers.CreateOvertimeHandler)
	r.Run(":8080")
}

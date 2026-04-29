
package main

import (
	"log"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/handlers"
	"github.com/eyoba-bisru/overtime-backend/internal/middleware"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/gin-contrib/cors"
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

	r.Use(cors.Default())

	base := r.Group("/api/v1")

	base.Use(middleware.LoggerMiddleware())

	auth := base.Group("/auth")
	auth.POST("/login", handlers.LoginHandler)

	// Auth routes that require authentication
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/change-password", handlers.ChangePasswordHandler)
	}

	admin := base.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RBACMiddleware(models.Admin))
	{
		admin.POST("/users", handlers.AdminCreateUserHandler)
		admin.GET("/departments", handlers.AdminGetDepartmentsHandler)
		admin.GET("/users", handlers.AdminGetUsersHandler)
		admin.PATCH("/users/:id", handlers.AdminUpdateUserHandler)
		admin.PATCH("/users/:id/block", handlers.AdminBlockUserHandler)
		admin.PATCH("/users/:id/reset-password", handlers.AdminResetPasswordHandler)
		admin.DELETE("/users/:id", handlers.AdminDeleteUserHandler)
	}

	overtime := base.Group("/overtime")
	overtime.Use(middleware.AuthMiddleware())
	{
		overtime.POST("", middleware.RBACMiddleware(models.Applicant, models.Admin), handlers.CreateOvertimeHandler)
		overtime.PATCH("/:id", middleware.RBACMiddleware(models.Applicant, models.Admin), handlers.UpdateOvertimeHandler)
		overtime.GET("/my", middleware.RBACMiddleware(models.Applicant, models.Admin), handlers.GetMyOvertimesHandler)
		overtime.GET("/pending", middleware.RBACMiddleware(models.Checker, models.Admin), handlers.GetPendingOvertimesHandler)
		overtime.PATCH("/:id/check", middleware.RBACMiddleware(models.Checker, models.Admin), handlers.CheckOvertimeHandler)
		overtime.GET("/checked", middleware.RBACMiddleware(models.Approver, models.Admin), handlers.GetCheckedOvertimesHandler)
		overtime.PATCH("/:id/approve", middleware.RBACMiddleware(models.Approver, models.Admin), handlers.ApproveOvertimeHandler)
		overtime.PATCH("/:id/reject", middleware.RBACMiddleware(models.Checker, models.Approver, models.Admin), handlers.RejectOvertimeHandler)
		overtime.GET("/approved", middleware.RBACMiddleware(models.Finance, models.Admin), handlers.GetApprovedOvertimesHandler)
		overtime.GET("/:id", handlers.GetOvertimeByIDHandler)
		overtime.DELETE("/:id", handlers.DeleteOvertimeHandler)
	}
	r.Run(":8080")
}

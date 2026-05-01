package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var loginInput LoginInput
	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	user := models.User{
		Email:        loginInput.Email,
		PasswordHash: loginInput.PasswordHash,
	}

	userData, token, forcePasswordChange, err := services.LoginService(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	host := os.Getenv("HOST")
	isSecure := !strings.Contains(host, "localhost") && !strings.Contains(host, "127.0.0.1")
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 86400, "/", host, isSecure, true)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: gin.H{
			"force_password_change": forcePasswordChange,
			"user": gin.H{
				"id":    userData.ID,
				"email": userData.Email,
				"name":  userData.Name,
				"role":  userData.Role,
			},
		},
	})
}

type ChangePasswordInput struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func ChangePasswordHandler(c *gin.Context) {
	var input ChangePasswordInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)

	err := services.ChangePasswordService(user.ID, input.NewPassword, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Password changed successfully"})
}

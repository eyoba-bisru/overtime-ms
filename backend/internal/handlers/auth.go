package handlers

import (
	"net/http"
	"os"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type UserInput struct {
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password" binding:"required"`
	Name         string `json:"name" binding:"required"`
}

func CreateUserHandler(c *gin.Context) {
	var userInput UserInput
	if err := c.ShouldBind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:        userInput.Email,
		PasswordHash: userInput.PasswordHash,
		Name:         userInput.Name,
	}

	data, err := services.CreateUserService(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": data})

}

func LoginHandler(c *gin.Context) {
	var userInput UserInput
	if err := c.ShouldBind(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:        userInput.Email,
		PasswordHash: userInput.PasswordHash,
	}

	token, err := services.LoginService(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", os.Getenv("HOST"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

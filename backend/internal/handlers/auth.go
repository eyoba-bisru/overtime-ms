package handlers

import (
	"net/http"
	"os"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	var user models.User
	c.Bind(&user)
	data, err := services.CreateUserService(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": data})

}

func LoginHandler(c *gin.Context) {
	var user models.User
	c.Bind(&user)
	token, err := services.LoginService(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", os.Getenv("HOST"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

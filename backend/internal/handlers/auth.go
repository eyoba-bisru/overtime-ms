package handlers

import (
	"net/http"

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

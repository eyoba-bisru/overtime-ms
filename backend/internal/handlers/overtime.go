package handlers

import (
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateOvertimeHandler(c *gin.Context) {
	var input models.Overtime
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	input.UserID = c.MustGet("user").(*models.User).ID
	// Process the overtime request
	data, err := services.CreateOvertimeService(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Overtime request created successfully", "data": data})
}

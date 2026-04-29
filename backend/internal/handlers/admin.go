package handlers

import (
	"net/http"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminCreateUserInput struct {
	Email        string      `json:"email" binding:"required,email"`
	PasswordHash string      `json:"password" binding:"required"`
	Name         string      `json:"name" binding:"required"`
	Role         models.Role `json:"role" binding:"required"`
}

type AdminUpdateUserInput struct {
	Email string      `json:"email" binding:"required,email"`
	Name  string      `json:"name" binding:"required"`
	Role  models.Role `json:"role" binding:"required"`
}

func AdminCreateUserHandler(c *gin.Context) {
	var input AdminCreateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Name:         input.Name,
		Role:         input.Role,
	}

	data, err := services.CreateUserService(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{Success: true, Message: "User created successfully", Data: data})
}

func AdminGetUsersHandler(c *gin.Context) {
	users, err := services.GetUsersService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: users})
}

func AdminUpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	var input AdminUpdateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	err := services.AdminUpdateUserService(id, input.Email, input.Name, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "User updated successfully"})
}

type AdminBlockUserInput struct {
	IsBlocked bool `json:"is_blocked"`
}

func AdminBlockUserHandler(c *gin.Context) {
	id := c.Param("id")
	var input AdminBlockUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	err := services.BlockUserService(id, input.IsBlocked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "User block status updated successfully"})
}

func AdminResetPasswordHandler(c *gin.Context) {
	id := c.Param("id")

	tempPassword, err := services.ResetUserPasswordService(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User password reset successfully",
		Data:    gin.H{"temporary_password": tempPassword},
	})
}

func AdminDeleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteUserService(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "User deleted successfully"})
}

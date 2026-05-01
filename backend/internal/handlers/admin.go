package handlers

import (
	"net/http"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminCreateUserInput struct {
	Email        string      `json:"email" binding:"required,email"`
	PasswordHash string      `json:"password" binding:"required"`
	Name         string      `json:"name" binding:"required"`
	Role         models.Role `json:"role" binding:"required"`
	DepartmentID string      `json:"department_id" binding:"required"`
}

type AdminUpdateUserInput struct {
	Email        string      `json:"email" binding:"required,email"`
	Name         string      `json:"name" binding:"required"`
	Role         models.Role `json:"role" binding:"required"`
	DepartmentID string      `json:"department_id" binding:"required"`
}

func GetCurrentUser(c *gin.Context) *models.User {
	val, exists := c.Get("user")
	if !exists {
		return nil
	}
	return val.(*models.User)
}

func AdminCreateUserHandler(c *gin.Context) {
	actor := GetCurrentUser(c)
	var input AdminCreateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	deptID, _ := uuid.Parse(input.DepartmentID)
	user := models.User{
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Name:         input.Name,
		Role:         input.Role,
		DepartmentID: deptID,
	}

	data, err := services.CreateUserService(&user, actor.ID)
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

func AdminGetDepartmentsHandler(c *gin.Context) {
	depts, err := repository.GetDepartmentsRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: depts})
}

func AdminUpdateUserHandler(c *gin.Context) {
	actor := GetCurrentUser(c)
	id := c.Param("id")
	var input AdminUpdateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	err := services.AdminUpdateUserService(id, input.Email, input.Name, input.Role, input.DepartmentID, actor.ID)
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
	actor := GetCurrentUser(c)
	id := c.Param("id")
	var input AdminBlockUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	err := services.BlockUserService(id, input.IsBlocked, actor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "User block status updated successfully"})
}

func AdminResetPasswordHandler(c *gin.Context) {
	actor := GetCurrentUser(c)
	id := c.Param("id")

	tempPassword, err := services.ResetUserPasswordService(id, actor.ID)
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
	actor := GetCurrentUser(c)
	id := c.Param("id")

	err := services.DeleteUserService(id, actor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "User deleted successfully"})
}

type DepartmentInput struct {
	Name string `json:"name" binding:"required"`
}

func AdminCreateDepartmentHandler(c *gin.Context) {
	actor := GetCurrentUser(c)
	var input DepartmentInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	id, err := repository.CreateDepartmentRepo(input.Name, actor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{Success: true, Message: "Department created successfully", Data: gin.H{"id": id}})
}

func AdminUpdateDepartmentHandler(c *gin.Context) {
	actor := GetCurrentUser(c)
	id := c.Param("id")
	var input DepartmentInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	parsedID, _ := uuid.Parse(id)
	err := repository.UpdateDepartmentRepo(parsedID, input.Name, actor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Department updated successfully"})
}

func AdminDeleteDepartmentHandler(c *gin.Context) {
	actor := GetCurrentUser(c)
	id := c.Param("id")
	parsedID, _ := uuid.Parse(id)

	err := repository.DeleteDepartmentRepo(parsedID, actor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Department deleted successfully"})
}

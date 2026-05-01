package handlers

import (
	"net/http"
	"strconv"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

func mapServiceError(err error) int {
	switch err {
	case models.ErrUnauthorized:
		return http.StatusForbidden
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrInvalidTransition:
		return http.StatusConflict
	case models.ErrImmutable:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func CreateOvertimeHandler(c *gin.Context) {
	var input models.Overtime
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)
	input.UserID = user.ID
	input.DepartmentID = user.DepartmentID
	data, err := services.CreateOvertimeService(input, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{Success: true, Message: "Overtime request created successfully", Data: data})
}

func UpdateOvertimeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: "Invalid ID"})
		return
	}

	var input models.Overtime
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)
	err = services.UpdateOvertimeService(id, user.ID, input, user.ID)
	if err != nil {
		c.JSON(mapServiceError(err), models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Overtime request updated successfully"})
}

func GetOvertimeByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: "Invalid ID"})
		return
	}

	data, err := services.GetOvertimeByIDService(id)
	if err != nil {
		c.JSON(mapServiceError(err), models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: data})
}

func paginatedOvertimeList(c *gin.Context, role models.Role, status models.OvertimeStatus) {
	user := c.MustGet("user").(*models.User)
	page, pageSize := parsePagination(c)

	data, total, err := services.GetOvertimesByStatusService(user.ID, role, status, user.DepartmentID.String(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "V2-TEST-ERROR", Error: err.Error()})
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "V2-TEST",
		Data:    data,
		Meta: &models.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func GetMyOvertimesHandler(c *gin.Context) {
	paginatedOvertimeList(c, models.Applicant, "")
}

func GetPendingOvertimesHandler(c *gin.Context) {
	paginatedOvertimeList(c, models.Checker, models.OvertimePending)
}

func GetCheckedOvertimesHandler(c *gin.Context) {
	paginatedOvertimeList(c, models.Approver, models.OvertimeChecked)
}

func GetApprovedOvertimesHandler(c *gin.Context) {
	paginatedOvertimeList(c, models.Finance, models.OvertimeApproved)
}

func AdminGetAllOvertimesHandler(c *gin.Context) {
	paginatedOvertimeList(c, models.Admin, "")
}

func CheckOvertimeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: "Invalid ID"})
		return
	}

	user := c.MustGet("user").(*models.User)
	err = services.CheckOvertimeService(id, user.Role, user.DepartmentID.String(), user.ID)
	if err != nil {
		c.JSON(mapServiceError(err), models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Overtime request checked successfully"})
}

func ApproveOvertimeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: "Invalid ID"})
		return
	}

	user := c.MustGet("user").(*models.User)
	err = services.ApproveOvertimeService(id, user.Role, user.DepartmentID.String(), user.ID)
	if err != nil {
		c.JSON(mapServiceError(err), models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Overtime request approved successfully"})
}

func RejectOvertimeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: "Invalid ID"})
		return
	}

	user := c.MustGet("user").(*models.User)
	err = services.RejectOvertimeService(id, user.Role, user.DepartmentID.String(), user.ID)
	if err != nil {
		c.JSON(mapServiceError(err), models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Overtime request rejected successfully"})
}

func DeleteOvertimeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Error: "Invalid ID"})
		return
	}

	user := c.MustGet("user").(*models.User)
	err = services.DeleteOvertimeService(id, user.ID, user.Role, user.ID)
	if err != nil {
		c.JSON(mapServiceError(err), models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "Overtime request deleted successfully"})
}


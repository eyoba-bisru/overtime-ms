package services

import (
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/google/uuid"
)

func calculateDuration(startTimeStr, endTimeStr string) (float64, error) {
	startTime, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		startTime, err = time.Parse("15:04:05", startTimeStr)
		if err != nil {
			return 0, err
		}
	}
	endTime, err := time.Parse("15:04", endTimeStr)
	if err != nil {
		endTime, err = time.Parse("15:04:05", endTimeStr)
		if err != nil {
			return 0, err
		}
	}

	duration := endTime.Sub(startTime).Hours()
	if duration < 0 {
		duration += 24 // Handle overnight overtime
	}
	return duration, nil
}

func CreateOvertimeService(input models.Overtime, actorID uuid.UUID) (uuid.UUID, error) {
	input.Status = models.OvertimePending
	input.ID = uuid.New()

	// Validate date format
	_, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return uuid.Nil, err
	}

	duration, err := calculateDuration(input.StartTime, input.EndTime)
	if err != nil {
		return uuid.Nil, err
	}
	input.Duration = duration

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	return repository.CreateOvertimeRepo(&input, actorID)
}

func UpdateOvertimeService(id uuid.UUID, userID uuid.UUID, input models.Overtime, actorID uuid.UUID) error {
	overtime, err := repository.GetOvertimeByIDRepo(id)
	if err != nil {
		return err
	}

	// Only the owner can edit, and only when pending
	if overtime.UserID != userID {
		return models.ErrUnauthorized
	}
	if overtime.Status != models.OvertimePending {
		return models.ErrImmutable
	}

	// Update fields
	if input.Date != "" {
		_, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			return err
		}
		overtime.Date = input.Date
	}
	if input.StartTime != "" {
		overtime.StartTime = input.StartTime
	}
	if input.EndTime != "" {
		overtime.EndTime = input.EndTime
	}
	if input.JobDone != "" {
		overtime.JobDone = input.JobDone
	}
	if input.Program != "" {
		overtime.Program = input.Program
	}

	// Recalculate duration
	duration, err := calculateDuration(overtime.StartTime, overtime.EndTime)
	if err != nil {
		return err
	}
	overtime.Duration = duration

	return repository.UpdateOvertimeRepo(overtime, actorID)
}

func GetOvertimeByIDService(id uuid.UUID) (*models.Overtime, error) {
	return repository.GetOvertimeByIDRepo(id)
}

func GetOvertimesByStatusService(userID uuid.UUID, role models.Role, status models.OvertimeStatus, departmentID string, page, pageSize int) ([]models.Overtime, int64, error) {
	total, err := repository.CountOvertimesRepo(userID, role, status, departmentID)
	if err != nil {
		return nil, 0, err
	}

	overtimes, err := repository.GetOvertimesRepo(userID, role, status, departmentID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return overtimes, total, nil
}

func CheckOvertimeService(id uuid.UUID, userRole models.Role, userDeptID string, actorID uuid.UUID) error {
	overtime, err := repository.GetOvertimeByIDRepo(id)
	if err != nil {
		return err
	}

	if userRole != models.Checker && userRole != models.Admin {
		return models.ErrUnauthorized
	}

	if userRole == models.Checker && overtime.DepartmentID.String() != userDeptID {
		return models.ErrUnauthorized
	}

	if overtime.Status != models.OvertimePending {
		return models.ErrInvalidTransition
	}

	return repository.UpdateOvertimeStatusRepo(id, models.OvertimeChecked, actorID)
}

func ApproveOvertimeService(id uuid.UUID, userRole models.Role, userDeptID string, actorID uuid.UUID) error {
	overtime, err := repository.GetOvertimeByIDRepo(id)
	if err != nil {
		return err
	}

	if userRole != models.Approver && userRole != models.Admin {
		return models.ErrUnauthorized
	}

	if userRole == models.Approver && overtime.DepartmentID.String() != userDeptID {
		return models.ErrUnauthorized
	}

	if overtime.Status != models.OvertimeChecked {
		return models.ErrInvalidTransition
	}

	return repository.UpdateOvertimeStatusRepo(id, models.OvertimeApproved, actorID)
}

func RejectOvertimeService(id uuid.UUID, userRole models.Role, userDeptID string, actorID uuid.UUID) error {
	overtime, err := repository.GetOvertimeByIDRepo(id)
	if err != nil {
		return err
	}

	if userRole != models.Admin && overtime.DepartmentID.String() != userDeptID && userRole != models.Finance {
		if userRole == models.Checker || userRole == models.Approver {
			return models.ErrUnauthorized
		}
	}

	if overtime.Status == models.OvertimeApproved || overtime.Status == models.OvertimeRejected {
		return models.ErrImmutable
	}

	if (userRole == models.Checker && overtime.Status == models.OvertimePending) ||
		(userRole == models.Approver && overtime.Status == models.OvertimeChecked) ||
		userRole == models.Admin {
		return repository.UpdateOvertimeStatusRepo(id, models.OvertimeRejected, actorID)
	}

	return models.ErrUnauthorized
}

func DeleteOvertimeService(id uuid.UUID, userID uuid.UUID, userRole models.Role, actorID uuid.UUID) error {
	overtime, err := repository.GetOvertimeByIDRepo(id)
	if err != nil {
		return err
	}

	if userRole == models.Admin || (overtime.UserID == userID && overtime.Status == models.OvertimePending) {
		return repository.DeleteOvertimeRepo(id, actorID)
	}

	return models.ErrUnauthorized
}

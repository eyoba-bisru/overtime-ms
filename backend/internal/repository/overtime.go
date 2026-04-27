package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/google/uuid"
)

func CreateOvertimeRepo(overtime *models.Overtime) (uuid.UUID, error) {
	var data uuid.UUID

	err := config.DB.QueryRow(context.Background(), "INSERT INTO overtimes (user_id, date, start_time, end_time, job_done, status, program, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", overtime.UserID, overtime.Date, overtime.StartTime, overtime.EndTime, overtime.JobDone, overtime.Status, overtime.Program, overtime.CreatedAt, overtime.UpdatedAt).Scan(&data)
	if err != nil {
		return uuid.Nil, err
	}
	return data, nil
}

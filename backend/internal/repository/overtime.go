package repository

import (
	"context"
	"fmt"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/google/uuid"
)

func CreateOvertimeRepo(overtime *models.Overtime) (uuid.UUID, error) {
	var data uuid.UUID

	err := config.DB.QueryRow(context.Background(), "INSERT INTO overtimes (user_id, date, start_time, end_time, job_done, status, program, duration, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", overtime.UserID, overtime.Date, overtime.StartTime, overtime.EndTime, overtime.JobDone, overtime.Status, overtime.Program, overtime.Duration, overtime.CreatedAt, overtime.UpdatedAt).Scan(&data)
	if err != nil {
		return uuid.Nil, err
	}
	return data, nil
}

func GetOvertimeByIDRepo(id uuid.UUID) (*models.Overtime, error) {
	var overtime models.Overtime
	err := config.DB.QueryRow(context.Background(), "SELECT id, user_id, date::TEXT, start_time::TEXT, end_time::TEXT, job_done, status, program, duration, created_at, updated_at FROM overtimes WHERE id = $1 AND deleted_at IS NULL", id).Scan(&overtime.ID, &overtime.UserID, &overtime.Date, &overtime.StartTime, &overtime.EndTime, &overtime.JobDone, &overtime.Status, &overtime.Program, &overtime.Duration, &overtime.CreatedAt, &overtime.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &overtime, nil
}

func buildOvertimeQuery(role models.Role, status models.OvertimeStatus, userID uuid.UUID) (string, string, []interface{}) {
	var where string
	var args []interface{}
	argIdx := 1

	switch role {
	case models.Applicant:
		where = fmt.Sprintf("user_id = $%d AND deleted_at IS NULL", argIdx)
		args = append(args, userID)
		argIdx++
	case models.Checker, models.Approver, models.Finance:
		where = fmt.Sprintf("status = $%d AND deleted_at IS NULL", argIdx)
		args = append(args, status)
		argIdx++
	case models.Admin:
		where = "deleted_at IS NULL"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM overtimes WHERE %s", where)
	selectQuery := fmt.Sprintf("SELECT id, user_id, date::TEXT, start_time::TEXT, end_time::TEXT, job_done, status, program, duration, created_at, updated_at FROM overtimes WHERE %s ORDER BY created_at DESC", where)

	return selectQuery, countQuery, args
}

func CountOvertimesRepo(userID uuid.UUID, role models.Role, status models.OvertimeStatus) (int64, error) {
	_, countQuery, args := buildOvertimeQuery(role, status, userID)

	var total int64
	err := config.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetOvertimesRepo(userID uuid.UUID, role models.Role, status models.OvertimeStatus, page, pageSize int) ([]models.Overtime, error) {
	selectQuery, _, args := buildOvertimeQuery(role, status, userID)

	offset := (page - 1) * pageSize
	argIdx := len(args) + 1
	selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := config.DB.Query(context.Background(), selectQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overtimes []models.Overtime
	for rows.Next() {
		var overtime models.Overtime
		err := rows.Scan(&overtime.ID, &overtime.UserID, &overtime.Date, &overtime.StartTime, &overtime.EndTime, &overtime.JobDone, &overtime.Status, &overtime.Program, &overtime.Duration, &overtime.CreatedAt, &overtime.UpdatedAt)
		if err != nil {
			return nil, err
		}
		overtimes = append(overtimes, overtime)
	}

	return overtimes, nil
}

func UpdateOvertimeRepo(overtime *models.Overtime) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE overtimes SET date = $1, start_time = $2, end_time = $3, job_done = $4, program = $5, duration = $6, updated_at = NOW() WHERE id = $7 AND deleted_at IS NULL", overtime.Date, overtime.StartTime, overtime.EndTime, overtime.JobDone, overtime.Program, overtime.Duration, overtime.ID)
	return err
}

func UpdateOvertimeStatusRepo(id uuid.UUID, status models.OvertimeStatus) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE overtimes SET status = $1, updated_at = NOW() WHERE id = $2", status, id)
	return err
}

func DeleteOvertimeRepo(id uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE overtimes SET deleted_at = NOW() WHERE id = $1", id)
	return err
}


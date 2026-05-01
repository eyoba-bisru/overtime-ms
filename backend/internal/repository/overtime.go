package repository

import (
	"context"
	"fmt"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/google/uuid"
)

func CreateOvertimeRepo(overtime *models.Overtime, actorID uuid.UUID) (uuid.UUID, error) {
	var data uuid.UUID
	err := config.DB.QueryRow(context.Background(), "INSERT INTO overtimes (user_id, date, start_time, end_time, job_done, status, program, department_id, duration, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12) RETURNING id", overtime.UserID, overtime.Date, overtime.StartTime, overtime.EndTime, overtime.JobDone, overtime.Status, overtime.Program, overtime.DepartmentID, overtime.Duration, overtime.CreatedAt, overtime.UpdatedAt, actorID).Scan(&data)
	if err != nil {
		return uuid.Nil, err
	}
	return data, nil
}

func GetOvertimeByIDRepo(id uuid.UUID) (*models.Overtime, error) {
	var overtime models.Overtime
	err := config.DB.QueryRow(context.Background(), `
		SELECT o.id, o.user_id, u.name as user_name, o.department_id, COALESCE(d.name, 'Unknown') as department_name, o.date::TEXT, o.start_time::TEXT, o.end_time::TEXT, o.job_done, o.status, o.program, o.duration, o.created_at, o.updated_at, o.created_by, o.updated_by
		FROM overtimes o
		JOIN users u ON o.user_id = u.id
		LEFT JOIN departments d ON o.department_id = d.id
		WHERE o.id = $1 AND o.deleted_at IS NULL`, id).Scan(&overtime.ID, &overtime.UserID, &overtime.UserName, &overtime.DepartmentID, &overtime.DepartmentName, &overtime.Date, &overtime.StartTime, &overtime.EndTime, &overtime.JobDone, &overtime.Status, &overtime.Program, &overtime.Duration, &overtime.CreatedAt, &overtime.UpdatedAt, &overtime.CreatedBy, &overtime.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &overtime, nil
}

func buildOvertimeQuery(role models.Role, status models.OvertimeStatus, userID uuid.UUID, departmentID string) (string, string, []interface{}) {
	var where string
	var args []interface{}
	argIdx := 1

	switch role {
	case models.Applicant:
		where = fmt.Sprintf("o.user_id = $%d AND o.deleted_at IS NULL", argIdx)
		args = append(args, userID)
		argIdx++
	case models.Checker, models.Approver:
		where = fmt.Sprintf("o.status = $%d AND o.department_id = $%d AND o.deleted_at IS NULL", argIdx, argIdx+1)
		args = append(args, status, departmentID)
		argIdx += 2
	case models.Finance:
		where = fmt.Sprintf("o.status = $%d AND o.deleted_at IS NULL", argIdx)
		args = append(args, status)
		argIdx++
	case models.Admin:
		where = "o.deleted_at IS NULL"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM overtimes o WHERE %s", where)
	selectQuery := fmt.Sprintf(`
		SELECT o.id, o.user_id, u.name as user_name, o.department_id, COALESCE(d.name, 'Unknown') as department_name, o.date::TEXT, o.start_time::TEXT, o.end_time::TEXT, o.job_done, o.status, o.program, o.duration, o.created_at, o.updated_at, o.created_by, o.updated_by
		FROM overtimes o
		JOIN users u ON o.user_id = u.id
		LEFT JOIN departments d ON o.department_id = d.id
		WHERE %s
		ORDER BY o.created_at DESC`, where)

	return selectQuery, countQuery, args
}

func CountOvertimesRepo(userID uuid.UUID, role models.Role, status models.OvertimeStatus, departmentID string) (int64, error) {
	_, countQuery, args := buildOvertimeQuery(role, status, userID, departmentID)

	var total int64
	err := config.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetOvertimesRepo(userID uuid.UUID, role models.Role, status models.OvertimeStatus, departmentID string, page, pageSize int) ([]models.Overtime, error) {
	selectQuery, _, args := buildOvertimeQuery(role, status, userID, departmentID)

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
		err := rows.Scan(&overtime.ID, &overtime.UserID, &overtime.UserName, &overtime.DepartmentID, &overtime.DepartmentName, &overtime.Date, &overtime.StartTime, &overtime.EndTime, &overtime.JobDone, &overtime.Status, &overtime.Program, &overtime.Duration, &overtime.CreatedAt, &overtime.UpdatedAt, &overtime.CreatedBy, &overtime.UpdatedBy)
		if err != nil {
			return nil, err
		}
		overtimes = append(overtimes, overtime)
	}

	return overtimes, nil
}

func UpdateOvertimeRepo(overtime *models.Overtime, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE overtimes SET date = $1, start_time = $2, end_time = $3, job_done = $4, program = $5, duration = $6, updated_at = NOW(), updated_by = $7 WHERE id = $8 AND deleted_at IS NULL", overtime.Date, overtime.StartTime, overtime.EndTime, overtime.JobDone, overtime.Program, overtime.Duration, actorID, overtime.ID)
	return err
}

func UpdateOvertimeStatusRepo(id uuid.UUID, status models.OvertimeStatus, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE overtimes SET status = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3 AND deleted_at IS NULL", status, actorID, id)
	return err
}

func DeleteOvertimeRepo(id uuid.UUID, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE overtimes SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2", actorID, id)
	return err
}


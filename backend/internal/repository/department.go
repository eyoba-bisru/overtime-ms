package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/google/uuid"
)

func GetDepartmentsRepo() ([]models.Department, error) {
	rows, err := config.DB.Query(context.Background(), "SELECT id, name FROM departments WHERE deleted_at IS NULL ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var depts []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name)
		if err != nil {
			return nil, err
		}
		depts = append(depts, d)
	}
	return depts, nil
}

func CreateDepartmentRepo(name string, actorID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := config.DB.QueryRow(context.Background(), "INSERT INTO departments (name, created_by, updated_by) VALUES ($1, $2, $2) RETURNING id", name, actorID).Scan(&id)
	return id, err
}

func UpdateDepartmentRepo(id uuid.UUID, name string, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE departments SET name = $1, updated_by = $2 WHERE id = $3 AND deleted_at IS NULL", name, actorID, id)
	return err
}

func DeleteDepartmentRepo(id uuid.UUID, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE departments SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2", actorID, id)
	return err
}

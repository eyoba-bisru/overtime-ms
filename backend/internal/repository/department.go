package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/google/uuid"
)

func GetDepartmentsRepo() ([]models.Department, error) {
	rows, err := config.DB.Query(context.Background(), "SELECT id, name FROM departments ORDER BY name ASC")
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

func CreateDepartmentRepo(name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := config.DB.QueryRow(context.Background(), "INSERT INTO departments (name) VALUES ($1) RETURNING id", name).Scan(&id)
	return id, err
}

func UpdateDepartmentRepo(id uuid.UUID, name string) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE departments SET name = $1 WHERE id = $2", name, id)
	return err
}

func DeleteDepartmentRepo(id uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "DELETE FROM departments WHERE id = $1", id)
	return err
}

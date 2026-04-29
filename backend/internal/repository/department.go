package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
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

package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/google/uuid"
)

func CreateUserRepo(user *models.User, actorID uuid.UUID) (string, error) {
	var id string
	err := config.DB.QueryRow(context.Background(), "INSERT INTO users (id, email, name, password_hash, role, department_id, force_password_change, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10) RETURNING id", user.ID, user.Email, user.Name, user.PasswordHash, user.Role, user.DepartmentID, user.ForcePasswordChange, user.CreatedAt, user.UpdatedAt, actorID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetUserByEmailRepo(email string) (*models.User, error) {
	row := config.DB.QueryRow(context.Background(), `
		SELECT u.id, u.email, u.name, u.password_hash, u.role, u.department_id, u.is_blocked, u.force_password_change, u.created_at, u.updated_at, COALESCE(d.name, 'Unknown') as department_name, u.created_by, u.updated_by
		FROM users u
		LEFT JOIN departments d ON u.department_id = d.id
		WHERE u.email = $1 AND u.deleted_at IS NULL`, email)
	user := &models.User{Department: &models.Department{}}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.DepartmentID, &user.IsBlocked, &user.ForcePasswordChange, &user.CreatedAt, &user.UpdatedAt, &user.Department.Name, &user.CreatedBy, &user.UpdatedBy)
	if err != nil {
		return nil, err
	}
	user.Department.ID = user.DepartmentID
	return user, nil
}

func UpdateUserRepo(user *models.User, actorID uuid.UUID) (*models.User, error) {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET email = $1, name = $2, password_hash = $3, role = $4, department_id = $5, last_login_at = $6, is_blocked = $7, force_password_change = $8, updated_at = $9, updated_by = $10 WHERE id = $11 AND deleted_at IS NULL", user.Email, user.Name, user.PasswordHash, user.Role, user.DepartmentID, user.LastLoginAt, user.IsBlocked, user.ForcePasswordChange, user.UpdatedAt, actorID, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUsersRepo() ([]models.User, error) {
	rows, err := config.DB.Query(context.Background(), `
		SELECT u.id, u.email, u.name, u.role, u.department_id, u.is_blocked, u.force_password_change, u.created_at, u.updated_at, COALESCE(d.name, 'Unknown') as department_name, u.created_by, u.updated_by
		FROM users u
		LEFT JOIN departments d ON u.department_id = d.id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{Department: &models.Department{}}
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.DepartmentID, &user.IsBlocked, &user.ForcePasswordChange, &user.CreatedAt, &user.UpdatedAt, &user.Department.Name, &user.CreatedBy, &user.UpdatedBy)
		if err != nil {
			return nil, err
		}
		user.Department.ID = user.DepartmentID
		users = append(users, user)
	}
	return users, nil
}

func UpdateUserRoleRepo(id string, role models.Role, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET role = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3 AND deleted_at IS NULL", role, actorID, id)
	return err
}

func UpdateUserPasswordRepo(id string, passwordHash string, force bool, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET password_hash = $1, force_password_change = $2, updated_at = NOW(), updated_by = $3 WHERE id = $4 AND deleted_at IS NULL", passwordHash, force, actorID, id)
	return err
}

func UpdateUserBlockStatusRepo(id string, isBlocked bool, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET is_blocked = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3 AND deleted_at IS NULL", isBlocked, actorID, id)
	return err
}

func DeleteUserRepo(id string, actorID uuid.UUID) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET deleted_at = NOW(), deleted_by = $1, updated_at = NOW() WHERE id = $2", actorID, id)
	return err
}

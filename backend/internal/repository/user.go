package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
)

func CreateUserRepo(user *models.User) (string, error) {
	var id string

	err := config.DB.QueryRow(context.Background(), "INSERT INTO users (id, email, name, password_hash, role, department_id, force_password_change, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", user.ID, user.Email, user.Name, user.PasswordHash, user.Role, user.DepartmentID, user.ForcePasswordChange, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetUserByEmailRepo(email string) (*models.User, error) {
	row := config.DB.QueryRow(context.Background(), `
		SELECT u.id, u.email, u.name, u.password_hash, u.role, u.department_id, u.is_blocked, u.force_password_change, u.created_at, u.updated_at, d.name as department_name
		FROM users u
		LEFT JOIN departments d ON u.department_id = d.id
		WHERE u.email = $1 AND u.deleted_at IS NULL`, email)
	user := &models.User{Department: &models.Department{}}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.DepartmentID, &user.IsBlocked, &user.ForcePasswordChange, &user.CreatedAt, &user.UpdatedAt, &user.Department.Name)
	if err != nil {
		return nil, err
	}
	user.Department.ID = user.DepartmentID
	return user, nil
}

func UpdateUserRepo(user *models.User) (*models.User, error) {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET email = $1, name = $2, password_hash = $3, role = $4, department_id = $5, last_login_at = $6, is_blocked = $7, force_password_change = $8, updated_at = $9 WHERE id = $10", user.Email, user.Name, user.PasswordHash, user.Role, user.DepartmentID, user.LastLoginAt, user.IsBlocked, user.ForcePasswordChange, user.UpdatedAt, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUsersRepo() ([]models.User, error) {
	rows, err := config.DB.Query(context.Background(), `
		SELECT u.id, u.email, u.name, u.role, u.department_id, u.is_blocked, u.force_password_change, u.created_at, u.updated_at, d.name as department_name
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
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.DepartmentID, &user.IsBlocked, &user.ForcePasswordChange, &user.CreatedAt, &user.UpdatedAt, &user.Department.Name)
		if err != nil {
			return nil, err
		}
		user.Department.ID = user.DepartmentID
		users = append(users, user)
	}
	return users, nil
}

func UpdateUserRoleRepo(id string, role models.Role) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL", role, id)
	return err
}

func UpdateUserPasswordRepo(id string, passwordHash string, force bool) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET password_hash = $1, force_password_change = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL", passwordHash, force, id)
	return err
}

func UpdateUserBlockStatusRepo(id string, isBlocked bool) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET is_blocked = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL", isBlocked, id)
	return err
}

func DeleteUserRepo(id string) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1", id)
	return err
}

package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
)

func CreateUserRepo(user *models.User) (string, error) {
	var id string

	err := config.DB.QueryRow(context.Background(), "INSERT INTO users (id, email, name, password_hash, role, force_password_change, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", user.ID, user.Email, user.Name, user.PasswordHash, user.Role, user.ForcePasswordChange, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetUserByEmailRepo(email string) (*models.User, error) {
	row := config.DB.QueryRow(context.Background(), "SELECT id, email, name, password_hash, role, is_blocked, force_password_change, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL", email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.IsBlocked, &user.ForcePasswordChange, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserRepo(user *models.User) (*models.User, error) {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET email = $1, name = $2, password_hash = $3, role = $4, last_login_at = $5, is_blocked = $6, force_password_change = $7, updated_at = $8 WHERE id = $9", user.Email, user.Name, user.PasswordHash, user.Role, user.LastLoginAt, user.IsBlocked, user.ForcePasswordChange, user.UpdatedAt, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUsersRepo() ([]models.User, error) {
	rows, err := config.DB.Query(context.Background(), "SELECT id, email, name, role, is_blocked, force_password_change, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.IsBlocked, &user.ForcePasswordChange, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
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

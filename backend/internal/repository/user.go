package repository

import (
	"context"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
)

func CreateUserRepo(user *models.User) (string, error) {
	var id string

	err := config.DB.QueryRow(context.Background(), "INSERT INTO users (id, email, name, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", user.ID, user.Email, user.Name, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetUserByEmailRepo(email string) (*models.User, error) {
	row := config.DB.QueryRow(context.Background(), "SELECT id, email, name, password_hash, role, created_at, updated_at FROM users WHERE email = $1", email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserRepo(user *models.User) (*models.User, error) {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET email = $1, name = $2, password_hash = $3, role = $4, last_login_at = $5, is_blocked = $6, updated_at = $7 WHERE id = $8", user.Email, user.Name, user.PasswordHash, user.Role, user.LastLoginAt, user.IsBlocked, user.UpdatedAt, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

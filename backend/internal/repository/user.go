package repository

import (
	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
)

func CreateUserRepo(user *models.User) (string, error) {
	var id string

	err := config.DB.QueryRow("INSERT INTO users (id, email, name, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", user.ID, user.Email, user.Name, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetUserByEmailRepo(email string) (*models.User, error) {
	row := config.DB.QueryRow("SELECT id, email, name, password, role, created_at, updated_at FROM users WHERE email = $1", email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

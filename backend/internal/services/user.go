package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserService(user *models.User) (sql.Result, error) {

	if user.Email == "" {
		return nil, errors.New("email is required")
	}
	if user.Password == "" {
		return nil, errors.New("password is required")
	}
	if user.Name == "" {
		return nil, errors.New("name is required")
	}

	role := models.Applicant // Default role
	user.Role = role
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	data, err := repository.CreateUserRepo(user)
	if err != nil {
		return nil, err
	}

	return data, nil
}

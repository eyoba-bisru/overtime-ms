package services

import (
	"errors"
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/utils"
	"github.com/google/uuid"
)

func CreateUserService(user *models.User) (string, error) {

	if user.Email == "" {
		return "", errors.New("email is required")
	}
	if user.PasswordHash == "" {
		return "", errors.New("password is required")
	}
	if user.Name == "" {
		return "", errors.New("name is required")
	}

	user.Role = models.Applicant
	user.ID = uuid.New()

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return "", err
	}
	user.PasswordHash = string(hashedPassword)

	data, err := repository.CreateUserRepo(user)
	if err != nil {
		return "", err
	}

	return data, nil
}

func LoginService(user *models.User) (string, error) {
	existingUser, err := repository.GetUserByEmailRepo(user.Email)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(user.PasswordHash, existingUser.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	if existingUser.IsBlocked {
		return "", errors.New("user is blocked")
	}

	now := time.Now()
	existingUser.LastLoginAt = &now
	_, err = repository.UpdateUserRepo(existingUser)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateJWT(existingUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

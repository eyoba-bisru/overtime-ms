package services

import (
	"database/sql"
	"errors"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/utils"
	"github.com/google/uuid"
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

	user.Role = models.Applicant
	user.ID = uuid.New()

	hashedPassword, err := utils.HashPassword(user.Password)
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

func LoginService(user *models.User) (string, error) {
	existingUser, err := repository.GetUserByEmailRepo(user.Email)
	if err != nil {
		return "", err
	}

	if !utils.CheckPasswordHash(user.Password, existingUser.Password) {
		return "", errors.New("invalid credentials")
	}

	if existingUser.IsBlocked {
		return "", errors.New("user is blocked")
	}

	token, err := utils.GenerateJWT(existingUser)
	if err != nil {
		return "", err
	}

	return token, nil
}

package services

import (
	"errors"
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/utils"
	"github.com/google/uuid"
)

func CreateUserService(user *models.User, actorID uuid.UUID) (string, error) {
	if user.Role == "" {
		user.Role = models.Applicant
	}
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return "", err
	}
	user.PasswordHash = string(hashedPassword)

	return repository.CreateUserRepo(user, actorID)
}

func LoginService(user *models.User) (*models.User, string, bool, error) {
	existingUser, err := repository.GetUserByEmailRepo(user.Email)
	if err != nil {
		return nil, "", false, err
	}

	if !utils.CheckPasswordHash(user.PasswordHash, existingUser.PasswordHash) {
		return nil, "", false, errors.New("invalid credentials")
	}

	if existingUser.IsBlocked {
		return nil, "", false, errors.New("user is blocked")
	}

	now := time.Now()
	existingUser.LastLoginAt = &now
	// Use the user themselves as the actor for updating last login
	_, err = repository.UpdateUserRepo(existingUser, existingUser.ID)
	if err != nil {
		return nil, "", false, err
	}

	token, err := utils.GenerateJWT(existingUser)
	if err != nil {
		return nil, "", false, err
	}

	return existingUser, token, existingUser.ForcePasswordChange, nil
}

func ChangePasswordService(userID uuid.UUID, newPassword string, actorID uuid.UUID) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return repository.UpdateUserPasswordRepo(userID.String(), hashedPassword, false, actorID)
}

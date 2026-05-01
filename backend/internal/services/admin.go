package services

import (
	"math/rand"
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/utils"
	"github.com/google/uuid"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

func generateRandomPassword(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetUsersService() ([]models.User, error) {
	return repository.GetUsersRepo()
}

func BlockUserService(id string, isBlocked bool, actorID uuid.UUID) error {
	return repository.UpdateUserBlockStatusRepo(id, isBlocked, actorID)
}

func ResetUserPasswordService(id string, actorID uuid.UUID) (string, error) {
	tempPassword := generateRandomPassword(12)
	hashedPassword, err := utils.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}

	err = repository.UpdateUserPasswordRepo(id, string(hashedPassword), true, actorID)
	if err != nil {
		return "", err
	}

	return tempPassword, nil
}

func DeleteUserService(id string, actorID uuid.UUID) error {
	return repository.DeleteUserRepo(id, actorID)
}

func AdminUpdateUserService(id string, email, name string, role models.Role, departmentID string, actorID uuid.UUID) error {
	parsedID, _ := uuid.Parse(id)
	deptID, _ := uuid.Parse(departmentID)
	user := &models.User{
		Base:         models.Base{ID: parsedID},
		Email:        email,
		Name:         name,
		Role:         role,
		DepartmentID: deptID,
	}

	_, err := repository.UpdateUserRepo(user, actorID)
	return err
}

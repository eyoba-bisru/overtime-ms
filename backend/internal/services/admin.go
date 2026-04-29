package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/config"
	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/utils"
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

func ChangeUserRoleService(id string, role models.Role) error {
	return repository.UpdateUserRoleRepo(id, role)
}

func BlockUserService(id string, isBlocked bool) error {
	return repository.UpdateUserBlockStatusRepo(id, isBlocked)
}

func ResetUserPasswordService(id string) (string, error) {
	tempPassword := generateRandomPassword(12)
	hashedPassword, err := utils.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}

	err = repository.UpdateUserPasswordRepo(id, string(hashedPassword), true)
	if err != nil {
		return "", err
	}

	return tempPassword, nil
}

func DeleteUserService(id string) error {
	return repository.DeleteUserRepo(id)
}

func AdminUpdateUserService(id string, email, name string, role models.Role) error {
	_, err := config.DB.Exec(context.Background(), "UPDATE users SET email = $1, name = $2, role = $3, updated_at = NOW() WHERE id = $4 AND deleted_at IS NULL", email, name, role, id)
	return err
}

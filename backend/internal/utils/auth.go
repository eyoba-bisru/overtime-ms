package utils

import (
	"os"
	"time"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(user *models.User) (string, error) {

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"dept_id": user.DepartmentID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(token string) (*models.User, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		deptID, _ := uuid.Parse(claims["dept_id"].(string))
		user := &models.User{
			Role:         models.Role(claims["role"].(string)),
			Email:        claims["email"].(string),
			DepartmentID: deptID,
			Base: models.Base{
				ID: uuid.MustParse(claims["id"].(string)),
			},
		}

		return user, nil
	}

	return nil, jwt.ErrInvalidKey

}

func HashPassword(password string) (string, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

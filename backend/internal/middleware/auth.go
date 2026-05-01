package middleware

import (
	"net/http"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/eyoba-bisru/overtime-backend/internal/repository"
	"github.com/eyoba-bisru/overtime-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user *models.User
		// check if the token is valid (you can implement your own logic here)
		if user, err = utils.ValidateJWT(tokenString); user == nil || err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check if user is blocked
		isBlocked, err := repository.GetUserBlockStatusRepo(user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if isBlocked {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is blocked"})
			return
		}

		c.Set("user", user)

		c.Next()

	}
}

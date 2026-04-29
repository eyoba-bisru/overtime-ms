package middleware

import (
	"net/http"

	"github.com/eyoba-bisru/overtime-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		u := user.(*models.User)
		authorized := false
		for _, role := range roles {
			if u.Role == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You don't have enough permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

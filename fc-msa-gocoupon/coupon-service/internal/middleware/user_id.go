package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserIDHeader = "X-USER-ID"
	userIDKey    = "userID"
)

// UserIDMiddleware validates the presence of the X-USER-ID header in the request.
// If present, it injects the user ID into the Gin context for downstream handlers.
// Returns a 400 Bad Request error if the header is missing.
func UserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader(UserIDHeader)
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "X-USER-ID header is required",
			})
			return
		}

		// Inject userID into Gin context
		c.Set(userIDKey, userID)

		c.Next()
	}
}

// GetCurrentUserID retrieves the user ID from the Gin context.
// Returns an error if the user ID is not set or invalid.
func GetCurrentUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return "", errors.New("user ID not found in context")
	}

	idStr, ok := userID.(string)
	if !ok || idStr == "" {
		return "", errors.New("user ID is invalid in context")
	}

	return idStr, nil
}

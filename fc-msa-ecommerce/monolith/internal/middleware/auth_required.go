package middleware

import (
	"net/http"
	"strings"

	"ecommerce/internal/model"
	"ecommerce/pkg/helper"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ErrUnauthorized())
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, "Authorization header format must be Bearer {token}", nil))
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := helper.JWTVerify(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, "Invalid or expired token", err))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

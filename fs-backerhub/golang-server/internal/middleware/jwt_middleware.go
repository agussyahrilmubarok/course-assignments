package middleware

import (
	"errors"
	"strings"

	"example.com.backend/internal/service"
	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(jwtService service.IJwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ex := exception.NewUnauthorized("Unauthorized", errors.New("authorization header is required"))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid authorization header format"))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		tokenString := parts[1]

		userID, err := jwtService.Validate(c.Request.Context(), tokenString)
		if err != nil {
			ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid or expired token"))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}

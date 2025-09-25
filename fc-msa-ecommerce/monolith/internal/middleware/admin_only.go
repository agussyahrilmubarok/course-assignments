package middleware

import (
	"ecommerce/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrUnauthorized())
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, "Invalid role format", nil))
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, model.ErrForbidden())
			c.Abort()
			return
		}

		c.Next()
	}
}

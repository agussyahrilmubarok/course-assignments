package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

const UserIDKey = "user_id"

func UserIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userID := c.Request().Header.Get("X-USER-ID")
			if userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing X-USER-ID header",
				})
			}

			ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

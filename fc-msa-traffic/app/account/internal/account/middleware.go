package account

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type ICustomMiddleware interface {
	Auth() echo.MiddlewareFunc
	RateLimiterConfig() middleware.RateLimiterConfig
	XUserID() echo.MiddlewareFunc
}

type customMiddleware struct {
	service IService
}

func NewCustomMiddleware(service IService) ICustomMiddleware {
	return &customMiddleware{
		service: service,
	}
}

// Auth middleware to validate JWT token and extract user_id into context
func (m *customMiddleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing Authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid Authorization header format"})
			}

			tokenStr := parts[1]
			userID, err := m.service.ValidateJwt(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or expired token"})
			}

			// Store userID in context for use in handlers
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// RateLimiterConfig returns a rate limiter configuration using memory store
func (m *customMiddleware) RateLimiterConfig() middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(1),   // 1 request per second
				Burst:     3,               // up to 3 requests in a short burst
				ExpiresIn: 3 * time.Minute, // token bucket expires after 3 minutes
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// You can customize this as needed (e.g., use user ID instead of IP)
			return c.Request().RemoteAddr, nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, echo.Map{
				"error": err.Error(),
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, echo.Map{
				"error": "Too many requests",
			})
		},
	}
}

// XUserID middleware to extract X-USER-ID from headers and put it in context
func (m *customMiddleware) XUserID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Request().Header.Get("X-USER-ID")
			if userID != "" {
				c.Set("user_id", userID)
			}
			return next(c)
		}
	}
}

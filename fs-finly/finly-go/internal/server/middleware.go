package server

import (
	"finly-go/internal/sessions"
	"finly-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequireAuth(c *fiber.Ctx) error {
	profileID, err := sessions.GetProfile(c)
	if err != nil {
		logger.Log.Error("Error retrieving profile from session", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving session")
	}

	if profileID == "" {
		return c.Redirect("/login")
	}

	return c.Next()
}

func PreventAfterAuth(c *fiber.Ctx) error {
	profileID, err := sessions.GetProfile(c)
	if err != nil {
		logger.Log.Error("Error retrieving profile from session", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving session")
	}

	if profileID != "" {
		if c.Path() == "/login" || c.Path() == "/signup" {
			return c.Redirect("/dashboard")
		}
	}

	return c.Next()
}

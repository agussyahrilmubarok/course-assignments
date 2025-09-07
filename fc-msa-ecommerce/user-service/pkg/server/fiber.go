package server

import (
	"example.com/user-service/pkg/config"
	"github.com/gofiber/fiber/v2"
)

func NewFiber(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               cfg.App.Name + " v" + cfg.App.Version,
		IdleTimeout:           cfg.App.IdleTimeout,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	return app
}

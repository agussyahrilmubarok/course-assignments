package controller

import "github.com/gofiber/fiber/v2"

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (hc *HomeController) GetIndex(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Welcome to Finly!",
	}

	return c.Render("index", data, "layouts/main")
}

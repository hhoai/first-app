package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func HomePageController(c *fiber.Ctx) error {
	return c.Render("homepage", fiber.Map{
		"Title": "Home Page",
	}, "layouts/main")
}

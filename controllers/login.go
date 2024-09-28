package controllers

import "github.com/gofiber/fiber/v2"

func LoginController(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

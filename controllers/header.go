package controller

import "github.com/gofiber/fiber/v2"

func GetUserName(c *fiber.Ctx) string {
	user := GetSessionUser(c)
	return user.Name
}

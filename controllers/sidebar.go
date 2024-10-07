package controller

import (
	"first-app/database"
	"first-app/models"
	"first-app/share"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetSessionUser(c *fiber.Ctx) models.User {
	var user models.User
	var role models.Role

	sess, _ := share.Store.Get(c)
	username := sess.Get("username")

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil { // binary: phan biet chu hoa, chu thuong
		log.Println("username not found in session")
	}

	sess.Set("role", role)

	return user
}

func RoleUser(c *fiber.Ctx) string {
	user := GetSessionUser(c)

	return user.Type
}

func CheckRoleUser(s string, c *fiber.Ctx) bool {
	permissions := RoleUser(c)

	log.Println(permissions, s, strings.EqualFold(permissions, s))
	return strings.EqualFold(permissions, s)
}

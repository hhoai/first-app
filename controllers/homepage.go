package controller

import (
	"first-app/database"
	"first-app/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func HomePageController(c *fiber.Ctx) error {
	user := GetSessionUser(c)

	var courses []models.Course

	if err := database.DB.Table("courses").
		Joins("JOIN course_users ON courses.course_id = course_users.course_id").
		Joins("JOIN users ON users.user_id = course_users.user_id").
		Where("users.username = ?", user.Username).
		Where("users.type = 'instructor'").
		Select("courses.course_id AS CourseID, courses.course_title AS course_title, courses.course_description AS course_description").
		Find(&courses).
		Error; err != nil {
		log.Println(err)
	}

	log.Println("course: ", courses)
	return c.Render("homepage", fiber.Map{
		"Title":   "Home Page",
		"Ctx":     c,
		"Courses": courses,
	}, "layouts/main")
}

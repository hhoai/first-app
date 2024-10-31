package controller

import (
	"first-app/database"
	"first-app/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CourseInstructorController(c *fiber.Ctx) error {
	id := c.Params("id")
	var course models.Course
	user := GetSessionUser(c)

	if err := database.DB.Where("course_id = ?", id).First(&course).Error; err != nil {
		log.Println(err)
	}

	var lesson []models.Lesson

	if err := database.DB.Where("course_id = ?", course.CourseID).Find(&lesson).Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"Ctx":               c,
		"CourseID":          course.CourseID,
		"CourseTitle":       course.CourseTitle,
		"CourseDescription": course.CourseDescription,
		"InstructorName":    user.Name,
		"Lessons":           lesson,
	}
	return c.Render("courseInstructor", data, "layouts/main")
}

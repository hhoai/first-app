package controller

import (
	"first-app/database"
	"first-app/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateAssignmentPostController(c *fiber.Ctx) error {

	var assignment models.Assignment

	c.BodyParser(&assignment)

	lessonID, _ := strconv.Atoi(c.Params("lessonID"))
	user := GetSessionUser(c)

	assignment.LessonID = lessonID
	assignment.UserID = user.UserID

	date, err := time.Parse("2006-01-02T15:04", c.FormValue("due_date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid date format")
	}

	newAssignment := models.Assignment{
		AssignmentTitle:  assignment.AssignmentTitle,
		AssignmentBody:   assignment.AssignmentBody,
		UserID:           user.UserID,
		LessonID:         lessonID,
		DueDate:          date,
		AssignmentStatus: "Đã giao",
	}

	database.DB.Create(&newAssignment)
	var lesson models.Lesson
	var assignments []models.Assignment

	if err := database.DB.Where("lesson_id = ?", lessonID).Find(&assignments).Error; err != nil {
		log.Println(err)
	}

	database.DB.Where("lesson_id = ?", lessonID).First(&lesson)

	data := fiber.Map{
		"Ctx":         c,
		"LessonID":    lessonID,
		"LessonTitle": lesson.LessonTitle,
		"CourseID":    lesson.CourseID,
		"Assignments": assignments,
	}

	return c.Render("lesson/detail", data, "layouts/main")
}

func AssignmentDetailController(c *fiber.Ctx) error {
	assignmentID := c.Params("assignmentID")

	var assignment models.Assignment

	if err := database.DB.Where("assignment_id = ?", assignmentID).First(&assignment).Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"Ctx":              c,
		"AssignmentID":     assignment.AssignmentID,
		"AssignmentTitle":  assignment.AssignmentTitle,
		"AssignmentBody":   assignment.AssignmentBody,
		"AssignmentStatus": assignment.AssignmentStatus,
		"DueDate":          assignment.DueDate,
	}

	return c.Render("assignment/detail", data, "layouts/main")
}

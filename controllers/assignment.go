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
		"DueDate":          assignment.DueDate.Format("02/01/2006"),
	}

	return c.Render("assignment/detail", data, "layouts/main")
}

func AssignmentDeleteController(c *fiber.Ctx) error {
	assignmentID := c.Params("assignmentID")

	var assignment models.Assignment
	if err := database.DB.Where("assignment_id = ?", assignmentID).First(&assignment).Error; err != nil {
		log.Println(err)
	}

	if err := database.DB.Delete(&assignment).Error; err != nil {
		log.Println(err)
	}
	lessonID := c.Params("lessonID")

	// hien thi lai
	var lesson models.Lesson
	var assignments []models.Assignment

	if err := database.DB.Where("lesson_id = ?", lessonID).First(&lesson).Error; err != nil {
		log.Println(err)
	}

	if err := database.DB.Where("lesson_id = ?", lessonID).Find(&assignments).Error; err != nil {
		log.Println(err)
	}

	var assignmentDetail []models.AssignmentDetail
	for _, assignment := range assignments {
		var temp models.AssignmentDetail
		temp.AssignmentID = assignment.AssignmentID
		temp.LessonID = assignment.LessonID
		temp.AssignmentBody = assignment.AssignmentBody
		temp.AssignmentStatus = assignment.AssignmentStatus
		temp.AssignmentTitle = assignment.AssignmentTitle
		temp.UserID = assignment.UserID
		temp.DueDate = assignment.DueDate.Format("15:04 02/01/2006")
		temp.CreatedAt = assignment.CreatedAt.Format("02/01/2006")
		assignmentDetail = append(assignmentDetail, temp)
	}

	data := fiber.Map{
		"Ctx":               c,
		"CourseID":          lesson.CourseID,
		"LessonID":          lesson.LessonID,
		"LessonTitle":       lesson.LessonTitle,
		"StartAt":           lesson.StartAt.Format("2006-01-02"),
		"LessonDescription": lesson.LessonDescription,
		"Assignments":       assignmentDetail,
	}

	return c.Render("lesson/detail", data, "layouts/main")
}

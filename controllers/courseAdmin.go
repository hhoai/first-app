package controller

import (
	"first-app/database"
	"first-app/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CourseController(c *fiber.Ctx) error {
	var courses []models.CourseWithRowNumber

	if err := database.DB.Table("courses").
		Joins("JOIN course_users ON courses.course_id = course_users.course_id").
		Joins("JOIN users ON course_users.user_id = users.user_id").
		Select("ROW_NUMBER() OVER (ORDER BY courses.course_id) AS RowNumber, courses.course_id AS CourseID, courses.course_title AS course_title, courses.course_description AS course_description, users.name AS Instructor").
		Find(&courses).
		Error; err != nil {
		log.Println(err)
	}

	log.Println(courses)

	data := fiber.Map{
		"Ctx":     c,
		"Courses": courses,
	}
	return c.Render("course/admin", data, "layouts/main")
}

func CreateCourseController(c *fiber.Ctx) error {

	data := fiber.Map{
		"Ctx":   c,
		"Title": "Thêm khóa học",
	}
	return c.Render("course/create", data, "layouts/main")
}

func CreateCoursePostController(c *fiber.Ctx) error {
	var p models.Course
	var existingCourse models.Course

	if err := c.BodyParser(&p); err != nil {
		log.Println(err.Error())
		return err
	}

	resultFindByCourseTitle := database.DB.Where("course_title = ?", p.CourseTitle).First(&existingCourse)

	if resultFindByCourseTitle.Error != nil && resultFindByCourseTitle.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		return c.Status(fiber.StatusInternalServerError).JSON(resultFindByCourseTitle.Error.Error())
	}

	if resultFindByCourseTitle.Error == nil {
		errorsMessage := fiber.Map{
			"RoleNameError": "Đã tồn tại khóa học!",
		}
		// Nếu không có lỗi và tìm thấy bản ghi, trả về lỗi trùng lặp
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	user := GetSessionUser(c)

	newCourse := models.Course{
		CourseTitle:       p.CourseTitle,
		CourseDescription: p.CourseDescription,
		UserID:            user.UserID,
	}

	if err := database.DB.Create(&newCourse).Error; err != nil {
		log.Println(err)
	}

	newCourseUser := models.CourseUser{
		CourseID: newCourse.CourseID,
		UserID:   3,
	}

	if err := database.DB.Create(&newCourseUser).Error; err != nil {
		log.Println(err)
	}

	// data := fiber.Map{
	// 	"Ctx": c,
	// }

	return c.Redirect("/admin/course")
}

func UpdateCourseController(c *fiber.Ctx) error {
	id := c.Params("id")

	var course models.Course
	var instructors []models.User

	if err := database.DB.Where("course_id = ?", id).First(&course).Error; err != nil {
		log.Println(err)
	}

	if err := database.DB.Where("type = 'instructor'").Find(&instructors).Error; err != nil {
		log.Println(err)
	}

	var instructor models.CourseUser

	if err := database.DB.Where("course_id = ?", id).First(&instructor).Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"CourseID":          course.CourseID,
		"CourseTitle":       course.CourseTitle,
		"CourseDescription": course.CourseDescription,
		"Ctx":               c,
		"Instructors":       instructors,
		"InstructorID":      instructor.UserID,
		"Title":             "Cập nhật khóa học",
	}

	return c.Render("course/edit", data, "layouts/main")
}

func UpdateCoursePutController(c *fiber.Ctx) error {
	id := c.Params("id")

	var p models.Course
	var courseUser models.CourseUser
	var existingCourse models.Course

	if err := c.BodyParser(&p); err != nil {
		log.Println(err.Error())
		return err
	}

	if err := c.BodyParser(&courseUser); err != nil {
		return err
	}

	nameError := validateAddress(p.CourseTitle)

	errorsMessage := fiber.Map{
		"RoleNameError": nameError,
	}

	if validateAddress(p.CourseTitle) != "" || validateAddress(p.CourseDescription) != "" {
		return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	}

	resultFindByCourseTitle := database.DB.Where("course_title = ?", p.CourseTitle).First(&existingCourse)

	if resultFindByCourseTitle.Error != nil && resultFindByCourseTitle.Error != gorm.ErrRecordNotFound {
		// Xử lý lỗi nếu có lỗi ngoài lỗi không tìm thấy
		return c.Status(fiber.StatusInternalServerError).JSON(resultFindByCourseTitle.Error.Error())
	}

	// if resultFindByCourseTitle.Error == nil {
	// 	errorsMessage := fiber.Map{
	// 		"RoleNameError": "Đã tồn tại khóa học!",
	// 	}
	// 	// Nếu không có lỗi và tìm thấy bản ghi, trả về lỗi trùng lặp
	// 	return c.Status(fiber.StatusConflict).JSON(errorsMessage)
	// }

	var updateCourse models.Course

	if err := database.DB.Where("course_id = ?", id).First(&updateCourse).Error; err != nil {
		log.Println(err)
	}

	updateCourse.CourseTitle = p.CourseTitle
	updateCourse.CourseDescription = p.CourseDescription

	if err := database.DB.Save(&updateCourse).Error; err != nil {
		log.Println(err)
	}

	var updateCourseUser models.CourseUser

	if err := database.DB.Where("course_id = ?", id).First(&updateCourseUser).Error; err != nil {
		log.Println(err)
	}

	updateCourseUser.UserID = courseUser.UserID

	log.Println("update", courseUser.UserID, updateCourseUser)

	if err := database.DB.Save(&updateCourseUser).Error; err != nil {
		log.Println(err)
	}

	var courses []models.CourseWithRowNumber

	if err := database.DB.Table("courses").
		Joins("JOIN course_users ON courses.course_id = course_users.course_id").
		Joins("JOIN users ON users.user_id = course_users.user_id").
		Select("ROW_NUMBER() OVER (ORDER BY courses.course_id) AS RowNumber, courses.course_id AS CourseID, courses.course_title AS course_title, courses.course_description AS course_description, users.name AS Instructor").
		Find(&courses).
		Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"Ctx":     c,
		"Courses": courses,
	}
	return c.Render("course/admin", data, "layouts/main")
}

func DeleteCourseController(c *fiber.Ctx) error {
	id := c.Params("id")

	var deleteCourse models.Course

	if err := database.DB.Delete(&deleteCourse, "course_id = ?", id).Error; err != nil {
		log.Println(err)
	}

	var courses []models.CourseWithRowNumber

	if err := database.DB.Table("courses").
		Select("ROW_NUMBER() OVER (ORDER BY courses.course_id) AS RowNumber, courses.course_id AS CourseID, courses.course_title AS course_title, courses.course_description AS course_description").
		Find(&courses).
		Error; err != nil {
		log.Println(err)
	}

	// data := fiber.Map{
	// 	"Ctx":     c,
	// 	"Courses": courses,
	// }
	return c.Redirect("/admin/course")
}

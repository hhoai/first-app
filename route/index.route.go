package route

import (
	"github.com/gofiber/fiber/v2"

	controller "first-app/controllers"
	"first-app/middleware"
)

func RouteInit(app *fiber.App) {
	// api
	api := app.Group("/api")

	api.Get("/account/admin", controller.GetDataAdmin)

	api.Get("/account/instructor", controller.GetDataInstructor)

	api.Get("/account/student", controller.GetDataStudent)

	// auth
	app.Get("/", middleware.IsAuthenticated, controller.HomePageController)

	app.Get("/login", controller.LoginController)

	app.Post("/login", controller.LoginPostController)

	app.Get("/logout", controller.LogoutController)

	// admin --- qltk
	// admin := app.Group("/admin", middleware.IsAuthenticated)

	app.Get("admin/account-admin", controller.AdminAccountController)

	app.Get("/admin/account-admin/account", controller.CreateAdminAccountController)

	app.Post("/admin/account-admin/account", controller.CreateAdminAccountPostController)

	app.Delete("/admin/account-admin/account/:id", controller.DeleteAdminAccountController)

	app.Get("/admin/account-admin/account/:id", controller.UpdateAdminAccountController)

	app.Put("/admin/account-admin/account/:id", controller.UpdateAdminAccountPutController)

	// admin --- qltk
	// admin := app.Group("/admin", middleware.IsAuthenticated)

	app.Get("admin/account-instructor", controller.InstructorAccountController)

	app.Get("/admin/account-instructor/account", controller.CreateInstructorAccountController)

	app.Post("/admin/account-instructor/account", controller.CreateInstructorAccountPostController)

	app.Delete("/admin/account-instructor/account/:id", controller.DeleteInstructorAccountController)

	app.Get("/admin/account-instructor/account/:id", controller.UpdateInstructorAccountController)

	app.Put("/admin/account-instructor/account/:id", controller.UpdateInstructorAccountPutController)

	// admin
	// ----- qltk
	admin := app.Group("/admin", middleware.IsAuthenticated)

	admin.Get("/account-student", controller.StudentAccountController)

	admin.Get("/account-student/account", controller.CreateStudentAccountController)

	admin.Post("/account-student/account", controller.CreateStudentAccountPostController)

	admin.Delete("/account-student/account/:id", controller.DeleteStudentAccountController)

	admin.Get("/account-student/account/:id", controller.UpdateStudentAccountController)

	admin.Put("/account-student/account/:id", controller.UpdateStudentAccountPutController)

	// ----- ql quyền

	admin.Get("/account/createRole", controller.CreateRoleController).Name("createRole")

	admin.Post("/account/createRole", controller.CreateRolePostController)

	admin.Delete("/account/role/:id", controller.DeleteRoleController).Name("deleteRole")

	admin.Get("/account/role/:id", controller.UpdateRoleController).Name("UpdateRole")

	admin.Put("/account/role/:id", controller.UpdateRolePutController)

	admin.Get("/account/role", controller.RoleController)

	// ----- ql khóa học
	admin.Get("/course", controller.CourseController)

	admin.Get("/course/:id", controller.CourseInstructorController)

	admin.Get("/course/createCourse", controller.CreateCourseController)

	admin.Post("/course/createCourse", controller.CreateCoursePostController)

	admin.Get("/course/updateCourse/:id", controller.UpdateCourseController)

	admin.Put("/course/updateCourse/:id", controller.UpdateCoursePutController)

	admin.Get("/course/deleteCourse/:id", controller.DeleteCourseController)

	// lesson
	admin.Get("/course/:id/lesson", controller.LessonController)

	admin.Post("/course/:id/lesson", controller.CreateLessonPostController)

	admin.Get("/course/lesson/:lessonID/detail", controller.LessonDetailController)

	admin.Delete("/course/lesson/:lessonID/detail", controller.LessonDeleteController)

	// assignment
	admin.Post("/course/lesson/:lessonID/detail", controller.CreateAssignmentPostController)

	admin.Get("/course/lesson/assignment/:assignmentID/detail", controller.AssignmentDetailController)
}

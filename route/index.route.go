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

	// admin --- qltk
	// admin := app.Group("/admin", middleware.IsAuthenticated)

	app.Get("admin/account-student", controller.StudentAccountController)

	app.Get("/admin/account-student/account", controller.CreateStudentAccountController)

	app.Post("/admin/account-student/account", controller.CreateStudentAccountPostController)

	app.Delete("/admin/account-student/account/:id", controller.DeleteStudentAccountController)

	app.Get("/admin/account-student/account/:id", controller.UpdateStudentAccountController)

	app.Put("/admin/account-student/account/:id", controller.UpdateStudentAccountPutController)
}

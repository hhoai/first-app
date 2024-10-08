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

	// ----- ql quyen

	admin.Get("/account/createRole", controller.CreateRoleController).Name("createRole")

	admin.Post("/account/createRole", controller.CreateRolePostController)

	admin.Get("/account/deleteRole/:id", controller.DeleteRoleController).Name("deleteRole")

	admin.Get("/account/updateRole/:id", controller.UpdateRoleController).Name("UpdateRole")

	admin.Put("/account/updateRole/:id", controller.UpdateRolePutController)

	admin.Get("/account/role", controller.RoleController)
}

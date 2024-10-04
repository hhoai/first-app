package main

import (
	controller "first-app/controllers"
	"first-app/database"
	"first-app/migration"
	"first-app/route"
	"first-app/share"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func main() {
	database.DatabaseInit()
	migration.RunMigration()

	share.Store = session.New()

	engine := html.New("./views", ".html")
	engine.AddFuncMap(fiber.Map{
		// "isPermissionSelected": controller.IsPermissionSelected,
		"checkRoleUser": controller.CheckRoleUser,
	})

	app := fiber.New(fiber.Config{
		Views: engine, // Sử dụng template engine đã nạp
	})

	app.Static("/", "./public")

	// // initial route
	route.RouteInit(app)

	app.Listen(":8080")
}

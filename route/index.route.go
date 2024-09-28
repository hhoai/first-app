package route

import (
	"github.com/gofiber/fiber/v2"

	"first-app/controllers"
)

func RouteInit(app *fiber.App) {
	app.Get("/", controllers.HomePageController)
	app.Get("/login", controllers.LoginController)
}

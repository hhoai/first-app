package route

import (
	"github.com/gofiber/fiber/v2"

	"first-app/controller"
)

func RouteInit(app *fiber.App) {
	app.Get("/", controller.HomePageController)
}

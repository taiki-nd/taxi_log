package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/api/v1/controller"
)

func Routes(app *fiber.App) {
	app.Get("/api/v1/users", controller.UsersIndex)
	app.Get("/api/v1/users/:id", controller.UsersShow)
	app.Post("/api/v1/users", controller.UsersCreate)
}

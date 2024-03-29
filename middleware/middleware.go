package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/taiki-nd/taxi_log/config"
)

func Cors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Content-Type, Content-Length, Authorization",
		AllowOrigins: config.Config.Url,
	}))
}

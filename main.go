package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/crons"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/middleware"
	"github.com/taiki-nd/taxi_log/routes"
	"github.com/taiki-nd/taxi_log/utils"
)

func main() {
	// logの有効化
	utils.Logging(config.Config.LogFile)

	// db接続
	db.ConnectToDb()

	// fiber
	app := fiber.New()

	// middleware
	middleware.Cors(app)

	// routes
	routes.Routes(app)

	// cron
	cron := cron.New()
	crons.CronManager(cron)
	cron.Start()

	app.Listen(":5050")

}

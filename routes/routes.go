package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/taiki-nd/taxi_log/controller/api/v1"
)

func Routes(app *fiber.App) {

	// user
	app.Get("/api/v1/users", controller.UsersIndex)
	app.Get("/api/v1/users/:id", controller.UsersShow)
	app.Post("/api/v1/users", controller.UsersCreate)
	app.Put("/api/v1/users/:id", controller.UsersUpdate)
	app.Delete("/api/v1/users/:id", controller.UsersDelete)
	app.Get("/api/v1/user/get_user_form_uid", controller.GetUserFromUuid)

	// area
	app.Get("/api/v1/areas", controller.GetAreas)

	// follow
	app.Post("/api/v1/follow", controller.Follow)
	app.Get("/api/v1/followings", controller.Followings)
	app.Get("/api/v1/followers", controller.Followers)
	app.Delete("/api/v1/followings", controller.DeleteFollowing)
	app.Post("/api/v1/follow_permission", controller.FollowPermission)

	// record
	app.Get("/api/v1/records", controller.RecordsIndex)
	app.Get("/api/v1/records/:id", controller.RecordsShow)
	app.Post("/api/v1/records", controller.RecordsCreate)
	app.Put("/api/v1/records/:id", controller.RecordsUpdate)
	app.Delete("/api/v1/records/:id", controller.RecordsDelete)

	// detail
	app.Get("/api/v1/details", controller.DetailsIndex)
	app.Get("/api/v1/details/:id", controller.DetailsShow)
	app.Post("/api/v1/details", controller.DetailsCreate)
	app.Put("/api/v1/details/:id", controller.DetailsUpdate)
	app.Delete("/api/v1/details/:id", controller.DetailsDelete)

	// analysis
	app.Get("/api/v1/analysis/sales_sum", controller.AnalysisSalesSum)
	app.Get("/api/v1/analysis/sales", controller.AnalysisSales)
	app.Get("/api/v1/analysis/records", controller.GetRecords)
	app.Get("/api/v1/analysis/analysis", controller.AnalysisPage)

	// ranking
	app.Get("/api/v1/ranking", controller.Ranking)
}

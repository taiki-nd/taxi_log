package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/service"
)

/**
 * Ranking
 */
func Ranking(c *fiber.Ctx) error {
	ranking_data, err := service.GetRankingData(c)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(fiber.Map{
		"data": ranking_data,
	})
}

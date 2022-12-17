package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * Ranking
 */
func Ranking(c *fiber.Ctx) error {
	ranking_data, err := service.GetRankingData(c)
	if err != nil {
		log.Printf("db_err: %v", err)
		service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}
	return service.SuccessResponse(c, []string{"ranking_data_success"}, ranking_data, nil)
}

package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * GetAreas
 */
func GetAreas(c *fiber.Ctx) error {
	areas, err := service.GetAreasFromPrefecture(c)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}
	return service.SuccessResponse(c, []string{"get_areas_success"}, areas, nil)
}

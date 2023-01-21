package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/service"
)

func GetProducts(c *fiber.Ctx) error {

	products, err := service.GetProducts()
	if err != nil {
		service.ErrorResponse(c, []string{"get_products_error"}, "failed to get products from stripe")
	}

	return service.SuccessResponse(c, nil, products, nil)
}

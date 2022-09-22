package service

import "github.com/gofiber/fiber/v2"

/**
 * ErrorResponse
 * error時のレスポンス内容
 * @params c *fiber.Ctx
 * @params code string
 * @params message string
 */
func ErrorResponse(c *fiber.Ctx, code interface{}, message string) error {
	return c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  false,
			"code":    code,
			"message": message,
		},
		"data": fiber.Map{},
	})
}

/**
 * SuccessResponse
 * error時のレスポンス内容
 * @params c *fiber.Ctx
 * @params code string
 * @params message string
 */
func SuccessResponse(c *fiber.Ctx, code string, data interface{}) error {
	return c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  true,
			"code":    code,
			"message": "",
		},
		"data": data,
	})
}

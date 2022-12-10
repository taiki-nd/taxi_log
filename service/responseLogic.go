package service

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

/**
 * ErrorResponse
 * error時のレスポンス内容
 * @params c *fiber.Ctx
 * @params code string
 * @params message string
 */
func ErrorResponse(c *fiber.Ctx, code []string, message string) error {
	log.Printf("error: %v", message)
	c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  false,
			"code":    code,
			"message": message,
		},
		"data": fiber.Map{},
	})
	return c.SendStatus(fiber.StatusBadRequest)
}

/**
 * SuccessResponse
 * error時のレスポンス内容
 * @params c *fiber.Ctx
 * @params code string
 * @params data interface{}
 * @params meta interface{}
 */
func SuccessResponse(c *fiber.Ctx, code []string, data interface{}, meta interface{}) error {
	c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  true,
			"code":    code,
			"message": "",
		},
		"data": data,
		"meta": meta,
	})
	return c.SendStatus(fiber.StatusOK)
}

/**
 * SuccessResponseAnalysis
 * error時のレスポンス内容
 * @params c *fiber.Ctx
 * @params code string
 * @params message string
 */
func SuccessResponseAnalysis(c *fiber.Ctx, code []string, data interface{}, meta interface{}) error {
	c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  true,
			"code":    code,
			"message": "",
		},
		"data":   data,
		"labels": meta,
	})
	return c.SendStatus(fiber.StatusOK)
}

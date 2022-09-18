package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/service"
)

func UsersIndex(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

/**
 * UsersCreate
 * userの新規登録処理
 * @params c *fiber.Ctx
 * @returns error error
 */
func UsersCreate(c *fiber.Ctx) error {
	// 変数確認
	var user *model.User

	// user認証

	// リクエストボディーのパース
	err := c.BodyParser(&user)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "body_parse_error",
				"message": fmt.Sprintf("body parse error: %v", err),
			},
			"data": fiber.Map{},
		})
	}

	// バリデーション
	_, errs := service.UserValidation(user)
	if len(errs) != 0 {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    errs,
				"message": fmt.Sprintf("validation error: %v", errs),
			},
			"data": fiber.Map{},
		})
	}

	// close_dayが31の場合の締め日の調整
	if user.CloseDay == 31 {
		date := service.AdjustmentCloseDay()
		user.CloseDay = date
	}

	// レコード作成
	err = db.DB.Create(&user).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "db_error",
				"message": fmt.Sprintf("db error: %v", err),
			},
			"data": fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  true,
			"code":    "create_user_success",
			"message": "",
		},
		"data": user,
	})
}

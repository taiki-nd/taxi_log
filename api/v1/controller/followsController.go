package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/service"
)

/**
 * Follow
 * フォロー機能
 * @params c *fiber.Ctx
 * @returns error
 */
func Follow(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.SuccessResponse(c, "user_auth_error", fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, "record_not_signin", "record not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, "user_not_match", "user not match")
	}

	var follow *model.Follow

	// リクエストボディーのパース
	err = c.BodyParser(&follow)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, "body_parse_error", fmt.Sprintf("body parse error: %v", err))
	}

	// レコード作成
	err = db.DB.Table("user_followings").Create(&follow).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "follow_user_success", follow)
}
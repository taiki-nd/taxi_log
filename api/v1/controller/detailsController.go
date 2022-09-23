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
 * DetailsIndex
 * detailsの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func DetailsIndex(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.ErrorResponse(c, "user_auth_error", fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, "user_not_signin", "user not signin")
	}

	// detailの検索
	details, err := service.SearchDetail(c, statuses[1])
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "index_detail_success", details)
}

/**
 * DetailsShow
 * detailの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func DetailsShow(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.ErrorResponse(c, "user_auth_error", fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, "user_not_signin", "user not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, "user_not_match", "user not match")
	}

	// レコードの取得
	detail, err := service.GetDetail(c)
	if err != nil {
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "show_detail_success", detail)
}

/**
 * DetailsCreate
 * detailの新規登録処理
 * @params c *fiber.Ctx
 * @returns error error
 */
func DetailsCreate(c *fiber.Ctx) error {
	// 変数確認
	var detail *model.Detail

	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.ErrorResponse(c, "user_auth_error", fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, "user_not_signin", "user not signin")
	}

	// リクエストボディーのパース
	err = c.BodyParser(&detail)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, "body_parse_error", fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs = service.DetailValidation(detail)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// レコード作成
	err = db.DB.Create(&detail).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "create_detail_success", detail)
}

/**
 * DetailsUpdate
 * detail情報の更新処理
 * @params c *fiber.Ctx
 */
func DetailsUpdate(c *fiber.Ctx) error {
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
		return service.ErrorResponse(c, "detail_not_signin", "detail not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, "user_not_match", "user not match")
	}

	// detailレコードの取得
	detail, err := service.GetDetail(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// リクエストボディのパース
	err = c.BodyParser(detail)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, "body_parse_error", fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs = service.DetailValidation(detail)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// detail情報の更新
	err = db.DB.Model(&detail).Updates(detail).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "update_detail_success", detail)
}

/**
 * DetailsDelete
 * detail情報の削除
 * @params c *fiber.Ctx
 * @returns error
 */
func DetailsDelete(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.ErrorResponse(c, "user_auth_error", fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, "user_not_signin", "user not signin")
	}
	// admin権限の確認
	if !statuses[1] {
		// user合致確認
		if !statuses[2] {
			return service.ErrorResponse(c, "user_not_match", "user not match")
		}
	}

	// detailレコードの取得
	detail, err := service.GetDetail(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// detail情報の削除
	errDetail := db.DB.Delete(detail).Error
	if errDetail != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "create_delete_success", nil)
}

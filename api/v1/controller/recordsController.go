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
 * RecordsIndex
 * recordの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsIndex(c *fiber.Ctx) error {
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

	// recordの検索
	records, err := service.SearchRecord(c, statuses[1])
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "index_record_success", records)
}

/**
 * RecordsShow
 * recordの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsShow(c *fiber.Ctx) error {
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
	record, err := service.GetRecord(c)
	if err != nil {
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// admin権限の確認
	if !statuses[1] {
		// follower確認
		status, err := service.IsFollowerForRecord(c, record)
		if err != nil {
			return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
		}
		if !status {
			return service.ErrorResponse(c, "follow_relationship_error", "follow relationship error")
		}
	}

	return service.SuccessResponse(c, "show_record_success", record)
}

/**
 * RecordsCreate
 * recordの新規登録処理
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsCreate(c *fiber.Ctx) error {
	// 変数確認
	var record *model.Record

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
	err = c.BodyParser(&record)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, "body_parse_error", fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs = service.RecordValidation(record)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// レコード作成
	err = db.DB.Create(&record).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "create_record_success", record)
}

/**
 * RecordsUpdate
 * record情報の更新処理
 * @params c *fiber.Ctx
 */
func RecordsUpdate(c *fiber.Ctx) error {
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

	// recordレコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// リクエストボディのパース
	err = c.BodyParser(record)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, "body_parse_error", fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs = service.RecordValidation(record)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// record情報の更新
	err = db.DB.Model(&record).Updates(record).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "update_record_success", record)
}

/**
 * RecordsDelete
 * record情報の削除
 * @params c *fiber.Ctx
 * @returns error
 */
func RecordsDelete(c *fiber.Ctx) error {
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

	// recordレコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// record情報の削除
	errRecord := db.DB.Delete(record).Error
	if errRecord != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "create_delete_success", nil)
}

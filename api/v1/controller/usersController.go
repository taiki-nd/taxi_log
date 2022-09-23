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
 * UsersIndex
 * userの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func UsersIndex(c *fiber.Ctx) error {
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
		return service.ErrorResponse(c, "user_not_signin", "user not sign in")
	}

	// userの検索
	users, err := service.SearchUser(c, statuses[1])
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "index_user_success", users)
}

/**
 * UsersShow
 * userの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func UsersShow(c *fiber.Ctx) error {
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

	// レコードの取得
	user, err := service.GetUser(c)
	if err != nil {
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "show_user_success", user)
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
		return service.ErrorResponse(c, "body_parse_error", fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs := service.UserValidation(user)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
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
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "create_user_success", user)
}

/**
 * UsersUpdate
 * user情報の更新処理
 * @params c *fiber.Ctx
 */
func UsersUpdate(c *fiber.Ctx) error {
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

	// userレコードの取得
	user, err := service.GetUser(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// リクエストボディのパース
	err = c.BodyParser(user)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// バリデーション
	_, errs = service.UserValidation(user)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// user情報の更新
	err = db.DB.Model(&user).Updates(user).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, "update_user_success", user)
}

/**
 * UsersDelete
 * user情報の削除
 * @params c *fiber.Ctx
 * @returns error
 */
func UsersDelete(c *fiber.Ctx) error {
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

	// userレコードの取得
	user, err := service.GetUser(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, "db_error", fmt.Sprintf("db error: %v", err))
	}

	// user削除トランザクション開始
	tx := db.DB.Begin()
	err = service.UsersDeleteTransaction(tx, user)
	if err != nil {
		tx.Rollback()
		log.Printf("transaction error: %v", err)
		return service.ErrorResponse(c, "transaction_error", fmt.Sprintf("transaction error: %v", err))
	}
	tx.Commit()

	return service.SuccessResponse(c, "delete_user_success", nil)
}

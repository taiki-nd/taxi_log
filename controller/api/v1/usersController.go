package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
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
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "user not sign in")
	}

	// userの検索
	users, err := service.SearchUser(c, statuses[1])
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"index_user_success"}, users, nil)
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
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "user not signin")
	}

	// admin権限の確認
	/*
		if !statuses[1] {
			// follower確認
			status, err := service.IsFollower(c)
			if err != nil {
				return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
			}
			if !status {
				return service.ErrorResponse(c, []string{constants.FOLLOW_RELATIONSHIP_ERROR}, "follow relationship error")
			}
		}
	*/

	// レコードの取得
	user, err := service.GetUser(c)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"show_user_success"}, user, nil)
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
		return service.ErrorResponse(c, []string{constants.BODY_PARSE_ERROR}, fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs := service.UserValidation(user)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// レコード作成
	err = db.DB.Create(&user).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"create_user_success"}, user, nil)
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
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "user not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	// userレコードの取得
	user, err := service.GetUser(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// リクエストボディのパース
	err = c.BodyParser(user)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// バリデーション
	_, errs = service.UserValidation(user)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	update_user := map[string]interface{}{
		"id":                user.Id,
		"uuid":              user.Uuid,
		"nickname":          user.Nickname,
		"profile_image_url": user.ProfileImageUrl,
		"prefecture":        user.Prefecture,
		"area":              user.Area,
		"company":           user.Company,
		"style_flg":         user.StyleFlg,
		"close_day":         user.CloseDay,
		"daily_target":      user.DailyTarget,
		"monthly_target":    user.MonthlyTarget,
		"is_tax":            user.IsTax,
		"open_flg":          user.OpenFlg,
		"is_admin":          user.IsAdmin,
	}

	// user情報の更新
	err = db.DB.Model(&user).Updates(update_user).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"update_user_success"}, user, nil)
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
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "user not signin")
	}
	// admin権限の確認
	if !statuses[1] {
		// user合致確認
		if !statuses[2] {
			return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
		}
	}

	// userレコードの取得
	user, err := service.GetUser(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// user削除トランザクション開始
	tx := db.DB.Begin()
	err = service.UsersDeleteTransaction(tx, user)
	if err != nil {
		tx.Rollback()
		log.Printf("transaction error: %v", err)
		return service.ErrorResponse(c, []string{"transaction_error"}, fmt.Sprintf("transaction error: %v", err))
	}
	tx.Commit()

	return service.SuccessResponse(c, []string{"delete_user_success"}, nil, nil)
}

/**
 * GetUserFromUuid
 * uuidからuser情報の取得
 * @params c *fiber.Ctx
 */
func GetUserFromUuid(c *fiber.Ctx) error {
	user, err := service.GetUserFromUuid(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"get_user_form_uuid_success"}, user, nil)
}

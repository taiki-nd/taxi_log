package controller

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * Followings
 * フォロー一覧の表示
 */
func Followings(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.SuccessResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "record not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	// 変数確認
	user_id_str := c.Query("user_id")
	user_id, _ := strconv.Atoi(user_id_str)

	// フォロイングのid一覧を取得する
	var following_ids []uint
	err = db.DB.Table("user_followings").Where("user_id = ?", uint(user_id)).Pluck("following_id", &following_ids).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// followingsの情報の取得
	var followings []*service.Follow
	err = db.DB.Table("users").Where("id IN (?)", following_ids).Find(&followings).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"show_followings_success"}, followings)
}

/**
 * Followers
 * フォロワー一覧の取得
 * @params c *fiber.Ctx
 * @returns error
 */
func Followers(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.SuccessResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "record not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	// 変数確認
	following_id_str := c.Query("following_id")
	following_id, _ := strconv.Atoi(following_id_str)

	// user_id一覧を取得
	var user_ids []uint
	err = db.DB.Table("user_followings").Where("following_id = ?", uint(following_id)).Pluck("user_id", &user_ids).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// usersの取得
	var users []*service.Follow
	err = db.DB.Table("users").Where("id IN (?)", user_ids).Find(&users).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"show_followers_success"}, users)
}

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
		return service.SuccessResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "record not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	var follow *model.Follow

	// リクエストボディーのパース
	err = c.BodyParser(&follow)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, []string{constants.BODY_PARSE_ERROR}, fmt.Sprintf("body parse error: %v", err))
	}

	// レコード作成
	err = db.DB.Table("user_followings").Create(&follow).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"follow_user_success"}, follow)
}

/*
 * DeleteFollowing
 * フォロー解除
 * @params c *fiber.Ctx
 * @returns error
 */
func DeleteFollowing(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.SuccessResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "record not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	// 変数確認
	user_id_str := c.Query("user_id")
	user_id, _ := strconv.Atoi(user_id_str)
	following_id_str := c.Query("following_id")
	following_id, _ := strconv.Atoi(following_id_str)

	// レコードの削除
	err = db.DB.Table("user_followings").Where("following_id = ?", uint(following_id)).Where("user_id = ?", uint(user_id)).Delete("").Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db_error: %v", err))
	}

	return service.SuccessResponse(c, []string{"delete_following_success"}, nil)
}

/**
 * FollowPermission
 * フォロー承認機能
 * @params c *fiber.Ctx
 * @returns error
 */
func FollowPermission(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.SuccessResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "record not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	// 変数確認
	var user_following model.UserFollowing
	user_id_str := c.Query("user_id")
	user_id, _ := strconv.Atoi(user_id_str)
	following_id_str := c.Query("following_id")
	following_id, _ := strconv.Atoi(following_id_str)

	// 対象レコードの取得
	err = db.DB.Where("following_id = ?", uint(user_id)).Where("user_id = ?", uint(following_id)).First(&user_following).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db_error: %v", err))
	}

	// permissionの付与
	update_user_following := map[string]interface{}{
		"id":           user_following.Id,
		"user_id":      user_following.UserId,
		"following_id": user_following.FollowingId,
		"permission":   !user_following.Permission,
	}

	// レコードの更新
	err = db.DB.Model(&user_following).Updates(update_user_following).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db_error: %v", err))
	}

	return service.SuccessResponse(c, []string{"permission_update_success"}, user_following)
}

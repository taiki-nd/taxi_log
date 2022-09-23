package service

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
)

type AuthUser struct {
	Id      uint
	Uuid    string
	IsAdmin bool `json:"is_admin"`
}

/**
 * UserAuth
 * user認証
 * @params c *fiber.Ctx
 * @returns bool
 * @returns string
 * @returns err
 */
func UserAuth(c *fiber.Ctx) ([]bool, []string, error) {
	// 変数の確認
	var statuses []bool
	var errs []string
	var authUser *AuthUser

	// headerの確認
	var header AuthUser
	err := c.ReqHeaderParser(&header)
	if err != nil {
		log.Println("reqHeader parse error")
		return statuses, errs, err
	}
	uuid := header.Uuid

	// user情報の取得
	if len(uuid) != 0 {
		err = db.DB.Table("users").Where("uuid = ?", uuid).First(&authUser).Error
		if err != nil {
			log.Printf("db_error: %v", err)
			return statuses, errs, err
		}
	}

	// signin状態の確認
	signin_status, err_signin := SignInCheck(c)
	if err_signin != "" {
		log.Printf("signin check error: %v", err_signin)
		statuses = append(statuses, signin_status)
		errs = append(errs, err_signin)
	}
	statuses = append(statuses, signin_status)

	if authUser != nil {
		// admin権限の確認
		statuses = append(statuses, authUser.IsAdmin)

		// user合致確認
		match_status, err_match := UserMatchCheck(c, authUser)
		log.Printf("match_status: %v", match_status)
		if err_match != "" {
			log.Printf("user match check error: %v", err_match)
			statuses = append(statuses, match_status)
			errs = append(errs, err_match)
		}
		statuses = append(statuses, match_status)
	} else {
		log.Printf("signin check error: %v", err_signin)
		statuses = append(statuses, signin_status)
		errs = append(errs, err_signin)
	}

	return statuses, errs, nil
}

/**
 * SignInCheck
 * signin状態の確認
 * @params c *fiber.Ctx
 * @returns bool
 * @returns string
 */
func SignInCheck(c *fiber.Ctx) (bool, string) {
	// headerの確認
	var header AuthUser
	err := c.ReqHeaderParser(&header)
	if err != nil {
		log.Println("reqHeader parse error")
		return false, "reqHeader_parse_error"
	}

	uuid := header.Uuid

	if len(uuid) == 0 {
		return false, ""
	}

	return true, ""
}

/**
 * UserMatchCheck
 * userの一致確認
 * @params user *model.User
 * @returns bool
 * @returns string
 */
func UserMatchCheck(c *fiber.Ctx, authUser *AuthUser) (bool, string) {
	// paramsの確認
	user_id := c.Query("user_id")
	user_id_int, _ := strconv.Atoi(user_id)

	if len(user_id) == 0 {
		return false, "check user match error"
	}

	// 合致確認
	if uint(user_id_int) == authUser.Id {
		return true, ""
	} else {
		return false, ""
	}
}

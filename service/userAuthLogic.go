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
	var user_id uint

	// headerの確認
	var header AuthUser
	err := c.ReqHeaderParser(&header)
	if err != nil {
		log.Println("reqHeader parse error")
		return statuses, errs, err
	}
	user_id = header.Id

	// SigninCheck
	if user_id == 0 {
		statuses = append(statuses, false)
		errs = append(errs, "user_not_signin")
	} else {
		statuses = append(statuses, true)
		errs = append(errs, "user_signin")
	}

	// AdminCheck
	err = db.DB.Table("users").Where("id = ?", user_id).First(&authUser).Error
	if err != nil {
		log.Printf("db_error: %v", err)
		return statuses, errs, err
	}
	if authUser.IsAdmin {
		statuses = append(statuses, true)
		errs = append(errs, "user_admin")
		statuses = append(statuses, false)
		errs = append(errs, "user_not_admin")
	} else {
		statuses = append(statuses, false)
		errs = append(errs, "user_not_admin")
	}

	//UserMatchCheck
	user_id_quey := c.Query("user_id")
	user_id_int_query, _ := strconv.Atoi(user_id_quey)
	if len(user_id_quey) == 0 {
		statuses = append(statuses, false)
		errs = append(errs, "user_not_match")
	} else {
		// 合致確認
		if uint(user_id_int_query) == authUser.Id {
			statuses = append(statuses, true)
			errs = append(errs, "user_match")
		} else {
			statuses = append(statuses, false)
			errs = append(errs, "user_not_match")
		}
	}

	return statuses, errs, nil
}

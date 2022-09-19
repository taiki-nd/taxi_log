package service

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
)

/**
 * SearchUser
 * usersの検索
 * @params c *fiber.Ctx
 * @returns users []*models.User
 */
func SearchUser(c *fiber.Ctx) ([]*model.User, error) {
	var users []*model.User

	// paramsの確認
	nickname := c.Query("nickname")
	prefecture := c.Query("prefecture")
	company := c.Query("company")
	style_flg := c.Query("style_flg")

	// クエリの作成
	userSearch := db.DB.Where("")
	if len(nickname) != 0 {
		userSearch.Where("nickname = ?", nickname)
	}
	if len(prefecture) != 0 {
		userSearch.Where("prefecture = ?", prefecture)
	}
	if len(company) != 0 {
		userSearch.Where("company = ?", company)
	}
	if len(style_flg) != 0 {
		userSearch.Where("style_flg = ?", style_flg)
	}
	// open_flgの確認
	// usersレコードの取得
	err := userSearch.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

/**
 * UserValidation
 * userのバリデーション機能
 * @params user *model.User
 * @returns bool
 * @returns []string
 */
func UserValidation(user *model.User) (bool, []string) {
	var errs []string

	// uuid
	if len(user.Uuid) == 0 {
		log.Println("uuid null error")
		errs = append(errs, "uuid_null_error")
	}
	//nickname
	if len(user.Nickname) == 0 {
		log.Println("nickname null error")
		errs = append(errs, "nickname_null_error")
	}
	if len(user.Nickname) < 3 || 30 < len(user.Nickname) {
		log.Println("nickname letter count error")
		errs = append(errs, "nickname_letter_count_error")
	}
	// prefecture
	if len(user.Prefecture) == 0 {
		log.Println("prefecture null error")
		errs = append(errs, "prefecture_null_error")
	}
	// company
	if len(user.Company) == 0 {
		log.Println("company null error")
		errs = append(errs, "company_null_error")
	}
	// style_flg
	if len(user.StyleFlg) == 0 {
		log.Println("style_flg null error")
		errs = append(errs, "style_flg_null_error")
	}
	if !(user.StyleFlg == "every_other_day" || user.StyleFlg == "day" || user.StyleFlg == "night" || user.StyleFlg == "other") {
		log.Println("specified word error(style_flg)")
		errs = append(errs, "specified_word_error(style_flg)")
	}
	// close_day
	if user.CloseDay == 0 {
		log.Println("close_day error")
		errs = append(errs, "close_day_error")
	}
	if user.CloseDay < 1 || 31 < user.CloseDay {
		log.Println("close_day date error")
		errs = append(errs, "close_day_date_error")
	}

	// errの出力
	if len(errs) != 0 {
		return false, errs
	}

	return true, nil
}

/**
 * AdjustmentCloseDay
 * 締め日の調整
 * @params time.Time
 * @returns time.Time
 */
func AdjustmentCloseDay() int64 {
	// 今
	now := time.Now()
	// 次月末
	close_day := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local)

	return int64(close_day.Day())
}

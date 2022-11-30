package service

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * SearchUser
 * usersの検索
 * @params c *fiber.Ctx
 * @returns users []*models.User
 */
func SearchUser(c *fiber.Ctx, adminStatus bool) ([]*model.User, error) {
	var users []*model.User
	log.Printf("adminStatus: %v", adminStatus)

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
	if !adminStatus {
		userSearch.Where("open_flg = ?", "open")
	}
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
	var searchedUser []*model.User

	// uuid検索
	err := db.DB.Table("users").Where("uuid = ?", user.Uuid).Find(&searchedUser).Error
	if err != nil {
		errs = append(errs, "db_error")
	}

	// uuid
	if len(user.Uuid) == 0 {
		log.Println("uuid null error")
		errs = append(errs, "uuid_null_error")
	}
	// uuid重複チェック
	if user.Id == 0 {
		if len(searchedUser) != 0 {
			errs = append(errs, "uuid_duplicate_error")
		}
	}
	//nickname
	if len(user.Nickname) == 0 {
		log.Println("nickname null error")
		errs = append(errs, "nickname_null_error")
	} else {
		if len(user.Nickname) < 3 || 30 < len(user.Nickname) {
			log.Println("nickname letter count error")
			errs = append(errs, "nickname_letter_count_error")
		}
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
	} else {
		if !(user.StyleFlg == "every_other_day" || user.StyleFlg == "day" || user.StyleFlg == "night" || user.StyleFlg == "other") {
			log.Println("specified word error(style_flg)")
			errs = append(errs, "specified_word_error(style_flg)")
		}
	}

	// close_day
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
 * @returns int64
 */
func AdjustmentCloseDay(year int, month int) int64 {
	close_day := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	return int64(close_day.Day())
}

/**
 * GetUser
 * user情報の取得
 * @params c *fiber.Ctx
 * @returns user *model.User
 */
func GetUser(c *fiber.Ctx) (*model.User, error) {
	// 変数確認
	user_id := c.Params("id")
	var user *model.User

	// レコードの取得
	err := db.DB.Where("id = ?", user_id).First(&user).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return nil, fmt.Errorf(constants.DB_ERR)
	}

	return user, nil
}

/**
 * GetUserFromUuid
 * uuidからuser情報を取得
 * @params c *fiber.Ctx
 * @returns *model.User
 */
func GetUserFromUuid(c *fiber.Ctx) (*model.User, error) {
	// headerの確認
	var header AuthUser
	err := c.ReqHeaderParser(&header)
	if err != nil {
		log.Println("reqHeader parse error")
		return nil, fmt.Errorf("reqHeader_parse_error")
	}
	uuid := header.Uuid

	var user *model.User
	err = db.DB.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return nil, fmt.Errorf(constants.DB_ERR)
	}

	return user, nil
}

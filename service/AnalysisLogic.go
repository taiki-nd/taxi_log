package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

func DataSettingForSalesSum(c *fiber.Ctx) ([]int64, error) {
	log.Println("DataSettingForSalesSum")

	// params
	//user_id, _ := strconv.Atoi(c.Query("user_id"))

	// 日報売上一覧の取得
	sales, err := GetSalesIndex(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return nil, fmt.Errorf(constants.DB_ERR)
	}

	return sales, nil
}

/**
 * GetSalesIndex
 * 日報売上一覧の取得
 * @params c *fiber.Ctx
 * @return []int64
 * @return error
 */
func GetSalesIndex(c *fiber.Ctx) ([]int64, error) {
	// params
	user_id, _ := strconv.Atoi(c.Query("user_id"))
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))

	// 締め日の取得
	user, err := GetUserFromUuid(c)
	if err != nil {
		return nil, err
	}
	close_day := user.CloseDay
	if user.CloseDay == 31 {
		date := AdjustmentCloseDay()
		close_day = date
	}
	sales_period_start := time.Date(year, time.Month(month), int(close_day), 12, 0, 0, 0, time.Local)
	if month >= 12 {
		year += 1
		month = 1
	} else {
		month += 1
	}
	sales_period_finish := time.Date(year, time.Month(month), int(close_day), 12, 0, 0, 0, time.Local)

	var sales []int64
	err = db.DB.Table("records").Where("user_id = ? && date > ? && date <= ?", user_id, sales_period_start, sales_period_finish).Order("date asc").Pluck("daily_sales", &sales).Error
	if err != nil {
		return nil, err
	}

	return sales, nil

}

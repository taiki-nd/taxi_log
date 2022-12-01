package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

func DataSettingForSalesSum(c *fiber.Ctx) ([]int64, []time.Time, error) {
	log.Println("DataSettingForSalesSum")

	// params
	//user_id, _ := strconv.Atoi(c.Query("user_id"))

	// 日報売上一覧の取得
	sales, dates, err := GetSalesIndex(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return nil, nil, fmt.Errorf(constants.DB_ERR)
	}

	var sales_sums []int64
	var i int64 = 0
	for _, daily_sales := range sales {
		i += daily_sales
		sales_sums = append(sales_sums, i)
	}

	return sales_sums, dates, nil
}

/**
 * GetSalesIndex
 * 日報売上一覧の取得
 * @params c *fiber.Ctx
 * @return []int64
 * @return error
 */
func GetSalesIndex(c *fiber.Ctx) ([]int64, []time.Time, error) {
	// params
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))
	var sales_period_start time.Time
	var sales_period_finish time.Time

	user, err := GetUserFromUuid(c)
	if err != nil {
		return nil, nil, err
	}

	// user_idの取得
	user_id := user.Id

	// 給与日の取得
	pay_day := user.PayDay

	// 締め日の取得
	close_day := user.CloseDay
	var close_day_start int64
	var close_day_finish int64

	// 取得するデータの期間の確定
	// 給与日が月をまたぐ場合
	if pay_day-close_day < 0 {
		// 開始期間の取得
		if month == 1 {
			year -= 1
			month = 11
		} else if month == 2 {
			year -= 1
			month = 12
		} else {
			month -= 2
		}
		if user.CloseDay == 31 {
			date := AdjustmentCloseDay(year, month)
			close_day_start = date
		} else {
			close_day_start = close_day
		}
		sales_period_start = time.Date(year, time.Month(month), int(close_day_start), 12, 0, 0, 0, time.Local)

		// 終了期間の取得
		if month == 12 {
			year += 1
			month = 1
		} else {
			month += 1
		}
		if user.CloseDay == 31 {
			date := AdjustmentCloseDay(year, month)
			close_day_finish = date
		} else {
			close_day_finish = close_day
		}
		sales_period_finish = time.Date(year, time.Month(month), int(close_day_finish), 12, 0, 0, 0, time.Local)

		// 給与日が月をまたがない場合
	} else {
		// 開始期間の取得
		if month == 1 {
			year -= 1
			month = 12
		} else {
			month -= 1
		}
		if user.CloseDay == 31 {
			date := AdjustmentCloseDay(year, month)
			close_day_start = date
		} else {
			close_day_start = close_day
		}
		sales_period_start = time.Date(year, time.Month(month), int(close_day_start), 12, 0, 0, 0, time.Local)

		// 終了期間の取得
		if month == 12 {
			year += 1
			month = 1
		} else {
			month += 1
		}
		if user.CloseDay == 31 {
			date := AdjustmentCloseDay(year, month)
			close_day_finish = date
		} else {
			close_day_finish = close_day
		}
		sales_period_finish = time.Date(year, time.Month(month), int(close_day_finish), 12, 0, 0, 0, time.Local)
	}

	var sales []int64
	err = db.DB.Table("records").Where("user_id = ? && date > ? && date <= ?", user_id, sales_period_start, sales_period_finish).Order("date asc").Pluck("daily_sales", &sales).Error
	if err != nil {
		return nil, nil, err
	}

	var dates []time.Time
	err = db.DB.Table("records").Where("user_id = ? && date > ? && date <= ?", user_id, sales_period_start, sales_period_finish).Order("date asc").Pluck("date", &dates).Error
	if err != nil {
		return nil, nil, err
	}

	return sales, dates, nil

}

/**
 * SearchRecordForMonth
 * recordsの検索
 * @params c *fiber.Ctx
 * @returns records []*models.Record
 */
func SearchRecordForMonth(c *fiber.Ctx) ([]*model.Record, error) {
	// params
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))

	user, err := GetUserFromUuid(c)
	if err != nil {
		return nil, err
	}

	// user_idの取得
	user_id := user.Id

	// 締め日の取得
	close_day := user.CloseDay
	if user.CloseDay == 31 {
		date := AdjustmentCloseDay(year, month)
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

	var records []*model.Record
	err = db.DB.Where("user_id = ? && date > ? && date <= ?", user_id, sales_period_start, sales_period_finish).Order("date asc").Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

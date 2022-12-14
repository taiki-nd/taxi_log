package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * SearchRecord
 * recordsの検索
 * @params c *fiber.Ctx
 * @returns records []*models.Record
 */
func SearchRecord(c *fiber.Ctx, adminStatus bool) ([]*model.Record, error) {
	var records []*model.Record

	// paramsの確認
	date := c.Query("date")
	day_of_week := c.Query("day_of_week")
	style_flg := c.Query("style_flg")
	occupancy_rate := c.Query("occupancy_rate")
	daily_sales := c.Query("daily_sales")

	// クエリの作成
	recordSearch := db.DB.Where("")
	if len(date) != 0 {
		recordSearch.Where("date = ?", date)
	}
	if len(day_of_week) != 0 {
		recordSearch.Where("day_of_week = ?", day_of_week)
	}
	if len(style_flg) != 0 {
		recordSearch.Where("style_flg = ?", style_flg)
	}
	if len(occupancy_rate) != 0 {
		recordSearch.Where("occupancy_rate >= ?", occupancy_rate)
	}
	if len(daily_sales) != 0 {
		recordSearch.Where("daily_sales >= ?", daily_sales)
	}
	// admin権限の確認
	user, err := GetUserFromQuery(c)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin {
		recordSearch.Where("user_id = ?", user.Id)
	}

	// pagination
	limit := constants.PAGINATION_LIMIT
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	// recordsレコードの取得
	err = recordSearch.Order("date desc").Offset(offset).Limit(limit).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

/**
 * RecordValidation
 * recordのバリデーション機能
 * @params record *model.Record
 * @returns bool
 * @returns []string
 */
func RecordValidation(record *model.Record) (bool, []string) {
	var errs []string

	// day_of_week
	if len(record.DayOfWeek) == 0 {
		log.Println("day_of_week null error")
		errs = append(errs, "day_of_week_null_error")
	}
	// style_flg
	if len(record.StyleFlg) == 0 {
		log.Println("style_flg null error")
		errs = append(errs, "style_flg_null_error")
	} else {
		if !(record.StyleFlg == "every_other_day" || record.StyleFlg == "day" || record.StyleFlg == "night" || record.StyleFlg == "other") {
			log.Println("specified word error(style_flg)")
			errs = append(errs, "specified_word_error(style_flg)")
		}
	}

	// start_hour
	if record.StartHour < 0 || 24 < record.StartHour {
		log.Println("start hour number error")
		errs = append(errs, "start_hour_number_error")
	}

	// running_time
	if record.RunningTime < 0 || 24 < record.RunningTime {
		log.Println("running time number error")
		errs = append(errs, "running_time_number_error")
	}

	// running_km
	if record.RunningKm < 0 {
		log.Println("running km number error")
		errs = append(errs, "running_km_number_error")
	}

	// number_of_time
	if record.NumberOfTime < 0 {
		log.Println("number of time number error")
		errs = append(errs, "number_of_time_number_error")
	}

	// occupancy_rate
	if record.OccupancyRate < 0 || 100 < record.OccupancyRate {
		log.Println("occupancy rate number error")
		errs = append(errs, "occupancy_rate_number_error")
	}

	// daily_sales
	if record.DailySales < 0 {
		log.Println("daily sales number error")
		errs = append(errs, "daily_sales_number_error")
	}

	// errの出力
	if len(errs) != 0 {
		return false, errs
	}

	return true, nil
}

/**
 * GetRecord
 * record情報の取得
 * @params c *fiber.Ctx
 * @returns record *model.Record
 */
func GetRecord(c *fiber.Ctx) (*model.Record, error) {
	// 変数確認
	record_id := c.Params("id")
	var record *model.Record

	// レコードの取得
	err := db.DB.Where("id = ?", record_id).First(&record).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return nil, fmt.Errorf(constants.DB_ERR)
	}

	return record, nil
}

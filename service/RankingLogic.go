package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/now"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * GetRankingData
 */
func GetRankingData(c *fiber.Ctx) (interface{}, error) {

	// 昨日の日付の取得
	today := now.BeginningOfDay()
	yesterday_start := today.AddDate(0, 0, -1)
	weekly_start := today.AddDate(0, 0, -7)
	monthly_start := today.AddDate(0, 0, -30)

	// 変数
	var daily_ranking_every_other_day_records []*model.Record
	var daily_ranking_day_records []*model.Record
	var daily_ranking_night_records []*model.Record

	var weekly_ranking_every_other_day_records []*model.Record
	var weekly_ranking_day_records []*model.Record
	var weekly_ranking_night_records []*model.Record

	var monthly_ranking_every_other_day_records []*model.Record
	var monthly_ranking_day_records []*model.Record
	var monthly_ranking_night_records []*model.Record

	// 昨日のランキング（隔日勤務）
	err := db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", yesterday_start, today, "every_other_day").Order("daily_sales DESC").Find(&daily_ranking_every_other_day_records).Error
	if err != nil {
		return nil, err
	}
	// 昨日のランキング（日勤）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", yesterday_start, today, "day").Order("daily_sales DESC").Find(&daily_ranking_day_records).Error
	if err != nil {
		return nil, err
	}
	// 昨日のランキング（夜勤）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", yesterday_start, today, "night").Order("daily_sales DESC").Find(&daily_ranking_night_records).Error
	if err != nil {
		return nil, err
	}

	// 7日間のランキング（隔日勤務）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", weekly_start, yesterday_start, "every_other_day").Order("daily_sales DESC").Find(&weekly_ranking_every_other_day_records).Error
	if err != nil {
		return nil, err
	}
	// 7日間のランキング（日勤）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", weekly_start, yesterday_start, "day").Order("daily_sales DESC").Find(&weekly_ranking_day_records).Error
	if err != nil {
		return nil, err
	}
	// 7日間のランキング（夜勤）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", weekly_start, yesterday_start, "night").Order("daily_sales DESC").Find(&weekly_ranking_night_records).Error
	if err != nil {
		return nil, err
	}

	// 31日間のランキング（隔日勤務）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", monthly_start, yesterday_start, "every_other_day").Order("daily_sales DESC").Find(&monthly_ranking_every_other_day_records).Error
	if err != nil {
		return nil, err
	}
	// 31日間のランキング（日勤）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", monthly_start, yesterday_start, "day").Order("daily_sales DESC").Find(&monthly_ranking_day_records).Error
	if err != nil {
		return nil, err
	}
	// 31日間のランキング（夜勤）
	err = db.DB.Limit(constants.RANKING_LIMIT).Table("records").Where("date >= ? && date <= ? && style_flg = ?", monthly_start, yesterday_start, "night").Order("daily_sales DESC").Find(&monthly_ranking_night_records).Error
	if err != nil {
		return nil, err
	}

	ranking_data := map[string]interface{}{
		"daily_ranking_every_other_day_records":   daily_ranking_every_other_day_records,
		"daily_ranking_day_records":               daily_ranking_day_records,
		"daily_ranking_night_records":             daily_ranking_night_records,
		"weekly_ranking_every_other_day_records":  weekly_ranking_every_other_day_records,
		"weekly_ranking_day_records":              weekly_ranking_day_records,
		"weekly_ranking_night_records":            weekly_ranking_night_records,
		"monthly_ranking_every_other_day_records": monthly_ranking_every_other_day_records,
		"monthly_ranking_day_records":             monthly_ranking_day_records,
		"monthly_ranking_night_records":           monthly_ranking_night_records,
	}

	return ranking_data, nil
}

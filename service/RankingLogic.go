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
	// ランキング表示エリアの取得
	prefecture := c.Query("prefecture")
	area := c.Query("area")

	rankingDataQuery := db.DB.Where("")

	if prefecture != "全エリア" {
		rankingDataQuery = db.DB.Where("prefecture = ?", prefecture)
	}

	if area != "全エリア" {
		rankingDataQuery = db.DB.Where("area = ?", area)
	}

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
	err := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", yesterday_start, today, "every_other_day").Order("daily_sales DESC").Find(&daily_ranking_every_other_day_records).Error
	if err != nil {
		return nil, err
	}
	// 昨日のランキング（日勤）
	err2 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", yesterday_start, today, "day").Order("daily_sales DESC").Find(&daily_ranking_day_records).Error
	if err2 != nil {
		return nil, err
	}
	// 昨日のランキング（夜勤）
	err3 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", yesterday_start, today, "night").Order("daily_sales DESC").Find(&daily_ranking_night_records).Error
	if err3 != nil {
		return nil, err
	}

	// 7日間のランキング（隔日勤務）
	err4 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", weekly_start, today, "every_other_day").Order("daily_sales DESC").Find(&weekly_ranking_every_other_day_records).Error
	if err4 != nil {
		return nil, err
	}
	// 7日間のランキング（日勤）
	err5 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", weekly_start, today, "day").Order("daily_sales DESC").Find(&weekly_ranking_day_records).Error
	if err5 != nil {
		return nil, err
	}
	// 7日間のランキング（夜勤）
	err6 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", weekly_start, today, "night").Order("daily_sales DESC").Find(&weekly_ranking_night_records).Error
	if err6 != nil {
		return nil, err
	}

	// 31日間のランキング（隔日勤務）
	err7 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", monthly_start, today, "every_other_day").Order("daily_sales DESC").Find(&monthly_ranking_every_other_day_records).Error
	if err7 != nil {
		return nil, err
	}
	// 31日間のランキング（日勤）
	err8 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", monthly_start, today, "day").Order("daily_sales DESC").Find(&monthly_ranking_day_records).Error
	if err8 != nil {
		return nil, err
	}
	// 31日間のランキング（夜勤）
	err9 := rankingDataQuery.Debug().Limit(constants.RANKING_LIMIT).Where("date >= ? && date <= ? && style_flg = ?", monthly_start, today, "night").Order("daily_sales DESC").Find(&monthly_ranking_night_records).Error
	if err9 != nil {
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

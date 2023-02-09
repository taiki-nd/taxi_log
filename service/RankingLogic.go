package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/now"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
)

// レスポンスモデルの定義
type Record struct {
	Id         uint      `json:"id"`
	Date       time.Time `json:"date"`
	StyleFlg   string    `json:"style_flg"`
	Prefecture string    `json:"prefecture"`
	Area       string    `json:"area"`
	DailySales int64     `json:"daily_sales"`
	UserId     uint      `json:"user_id"`
}

/**
 * GetRankingData
 */
func GetRankingData(c *fiber.Ctx) (interface{}, error) {
	// ランキング表示エリアの取得
	/*
		prefecture := c.Query("prefecture")
		area := c.Query("area")
	*/
	prefecture := "全エリア"
	area := "全エリア"

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

	//
	// 各期間のレコードを取得
	//

	// ランキング表示対象期間の全レコードを取得
	var records_for_month []Record
	err := rankingDataQuery.Model(model.Record{}).Where("date >= ? && date <= ?", monthly_start, today).Order("daily_sales DESC").Find(&records_for_month).Error
	if err != nil {
		return nil, err
	}

	// 1週間分のrecordを取得する
	var records_for_week []Record
	for _, record := range records_for_month {
		if today.After(record.Date) && record.Date.After(weekly_start) {
			records_for_week = append(records_for_week, record)
		}
	}

	// 1日分のrecordを取得する
	var records_for_day []Record
	for _, record := range records_for_month {
		if today.After(record.Date) && record.Date.After(yesterday_start) {
			records_for_day = append(records_for_day, record)
		}
	}

	//
	// ランキング情報の抽出
	//

	/*
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
	*/

	ranking_data := map[string]interface{}{
		/*
			"daily_ranking_every_other_day_records":   daily_ranking_every_other_day_records,
			"daily_ranking_day_records":               daily_ranking_day_records,
			"daily_ranking_night_records":             daily_ranking_night_records,
			"weekly_ranking_every_other_day_records":  weekly_ranking_every_other_day_records,
			"weekly_ranking_day_records":              weekly_ranking_day_records,
			"weekly_ranking_night_records":            weekly_ranking_night_records,
			"monthly_ranking_every_other_day_records": monthly_ranking_every_other_day_records,
			"monthly_ranking_day_records":             monthly_ranking_day_records,
			"monthly_ranking_night_records":           monthly_ranking_night_records,
		*/
	}

	return ranking_data, nil
}

func RecordsClassification(records []Record) ([]Record, []Record, []Record) {
	var records_every_other_day []Record
	var records_day []Record
	var records_night []Record

	for _, record := range records {
		if record.StyleFlg == "every_other_day" {
			records_every_other_day = append(records_every_other_day, record)
		} else if record.StyleFlg == "day" {
			records_day = append(records_day, record)
		} else if record.StyleFlg == "night" {
			records_night = append(records_night, record)
		} else {
			// nothing todo
		}
	}

	return records_every_other_day, records_day, records_night
}

type Records []Record

func RecordSort(records []Record) {

}

func (r Records) Len() int {
	return len(r)
}

func (r Records) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Records) Less(i, j int) bool {
	return r[i].DailySales >= r[j].DailySales
}

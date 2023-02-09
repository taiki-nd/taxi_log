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
	// 勤務形態別に分類
	//

	// 1カ月分の分類
	monthly_every_other_day_records, monthly_day_records, monthly_night_records := RecordsClassification(records_for_month)

	// 1週間分の分類
	weekly_every_other_day_records, weekly_day_records, weekly_night_records := RecordsClassification(records_for_week)

	// 1日分の分類
	daily_every_other_day_records, daily_day_records, daily_night_records := RecordsClassification(records_for_day)

	//
	// ランキング情報の抽出
	//

	var daily_ranking_every_other_day_records []Record
	var daily_ranking_day_records []Record
	var daily_ranking_night_records []Record

	var weekly_ranking_every_other_day_records []Record
	var weekly_ranking_day_records []Record
	var weekly_ranking_night_records []Record

	var monthly_ranking_every_other_day_records []Record
	var monthly_ranking_day_records []Record
	var monthly_ranking_night_records []Record

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

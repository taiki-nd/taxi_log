package service

import (
	"fmt"
	"sort"
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
	// TODO:暫定的に固定値で実行中
	prefecture := "全エリア"
	area := "全エリア"

	rankingDataQuery := db.DB.Where("")

	if prefecture != "全エリア" {
		rankingDataQuery = db.DB.Where("prefecture = ?", prefecture)
	}

	if area != "全エリア" {
		rankingDataQuery = db.DB.Where("area = ?", area)
	}

	// 日付の取得
	today := now.BeginningOfDay()              // 今日
	yesterday_start := today.AddDate(0, 0, -1) // 昨日
	weekly_start := today.AddDate(0, 0, -7)    // 1週前
	monthly_start := today.AddDate(0, 0, -30)  // 一カ月前 TODO:単純に30日減算で良いか？

	//
	// 各期間のレコードを取得
	//
	// ランキング表示対象期間の全レコードを取得
	var records_for_month []Record
	err := rankingDataQuery.Model(model.Record{}).
		Where("date >= ? && date <= ?", monthly_start, today).
		Order("daily_sales DESC").
		Find(&records_for_month).
		Error
	if err != nil {
		return nil, err
	}

	var records_for_week []Record
	var records_for_day []Record
	for _, record := range records_for_month {
		targetDate := record.Date
		// 1日分のrecordを取得する (昨日 <= 対象日X <= 今日)
		if targetDate.After(yesterday_start) && today.After(targetDate) {
			records_for_day = append(records_for_day, record)
		}
		// 1週間分のrecordを取得する (1週間前 <= 対処日X <= 今日)
		if targetDate.After(weekly_start) && today.After(targetDate) {
			records_for_week = append(records_for_week, record)
		}
	}

	//
	// 勤務形態別に分類
	//
	// 1カ月分の分類
	monthly_every_other_day_records, monthly_day_records, monthly_night_records := ClassifyRecords(records_for_month)

	// 1週間分の分類
	weekly_every_other_day_records, weekly_day_records, weekly_night_records := ClassifyRecords(records_for_week)

	// 1日分の分類
	daily_every_other_day_records, daily_day_records, daily_night_records := ClassifyRecords(records_for_day)

	// 並び替えてランキング情報として抽出
	ranking_data := map[string]interface{}{
		"daily_ranking_every_other_day_records":   SortRecordsByDailySales(daily_every_other_day_records),
		"daily_ranking_day_records":               SortRecordsByDailySales(daily_day_records),
		"daily_ranking_night_records":             SortRecordsByDailySales(daily_night_records),
		"weekly_ranking_every_other_day_records":  SortRecordsByDailySales(weekly_every_other_day_records),
		"weekly_ranking_day_records":              SortRecordsByDailySales(weekly_day_records),
		"weekly_ranking_night_records":            SortRecordsByDailySales(weekly_night_records),
		"monthly_ranking_every_other_day_records": SortRecordsByDailySales(monthly_every_other_day_records),
		"monthly_ranking_day_records":             SortRecordsByDailySales(monthly_day_records),
		"monthly_ranking_night_records":           SortRecordsByDailySales(monthly_night_records),
	}

	fmt.Printf("ranking_data: %v \n", ranking_data)

	return ranking_data, nil
}

func ClassifyRecords(records []Record) ([]Record, []Record, []Record) {
	var records_every_other_day []Record
	var records_day []Record
	var records_night []Record

	for _, record := range records {
		switch record.StyleFlg {
		case "every_other_day":
			records_every_other_day = append(records_every_other_day, record)
		case "day":
			records_day = append(records_day, record)
		case "night":
			records_night = append(records_night, record)
		}
	}

	return records_every_other_day, records_day, records_night
}

// TODO: 置き換え用関数 要テスト
func SortRecordsByDailySales_forReplace(records []Record) []Record {
	sortFn := func(i, j int) bool {
		return records[i].DailySales > records[j].DailySales
	}

	// イテレータを安定ソート
	sort.SliceStable(records, sortFn)
	if len(records) > 5 {
		return records[:5]
	}
	return records
}

// TODO: ループ過多. SortRecordsByDailySales_forReplace　に変更したい
func SortRecordsByDailySales(records []Record) []Record {
	var sortedRecords []Record

	// 売上のみ取得
	var dailySalesArray []int
	for _, record := range records {
		dailySalesArray = append(dailySalesArray, int(record.DailySales))
	}

	// 売上のソート
	sort.Sort(sort.Reverse(sort.IntSlice(dailySalesArray)))

	// recordのソート(上位5件のみ取得)
	for _, dailySales := range dailySalesArray {
		for _, record := range records {
			if dailySales == int(record.DailySales) {
				sortedRecords = append(sortedRecords, record)
				break
			}
		}
		if len(sortedRecords) == 5 {
			break
		}
	}

	return sortedRecords
}

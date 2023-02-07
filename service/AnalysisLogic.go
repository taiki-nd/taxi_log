package service

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

func DataSettingForSalesSum(c *fiber.Ctx) ([]int64, []time.Time, error) {

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

	user, err := GetUserFromQuery(c)
	if err != nil {
		return nil, nil, err
	}

	// userの取得
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
	var sales_period_start time.Time
	var sales_period_finish time.Time

	user, err := GetUserFromQuery(c)
	if err != nil {
		return nil, err
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

	var records []*model.Record
	err = db.DB.Where("user_id = ? && date > ? && date <= ?", user_id, sales_period_start, sales_period_finish).Order("date asc").Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

/**
 * GetAllAnalysisData
 * @params c *fiber.Ctx
 * @return
 */
func GetAllAnalysisData(c *fiber.Ctx) (interface{}, interface{}, interface{}, error) {
	user, err := GetUserFromQuery(c)
	if err != nil {
		return nil, nil, nil, err
	}

	var records []model.Record

	// params
	start_year, _ := strconv.Atoi(c.Query("start_year"))
	start_month, _ := strconv.Atoi(c.Query("start_month"))
	finish_year, _ := strconv.Atoi(c.Query("finish_year"))
	finish_month, _ := strconv.Atoi(c.Query("finish_month"))
	user_id := user.Id

	// 期間の設定
	period_start := time.Date(start_year, time.Month(start_month), 1, 0, 0, 0, 0, time.Local)
	period_finish := time.Date(finish_year, time.Month(finish_month+1), 1, 0, 0, 0, 0, time.Local)

	//
	// 対象期間データの取得
	//
	err = db.DB.Table("records").Where("user_id = ? && date > ? && date <= ?", user_id, period_start, period_finish).Find(&records).Error
	if err != nil {
		return nil, nil, nil, err
	}

	fmt.Println(records)

	// 曜日別売上解析データ
	average_sales_per_day := AnalysisAverageSalesPerDay(records)

	// 曜日別実車率データ
	average_occupancy_rate_per_day := AnalysisAverageOccupancyRatePerDay(records)

	// 期間内解析データの取得
	period_data := PeriodAnalysisData(records)

	return average_sales_per_day, average_occupancy_rate_per_day, period_data, nil
}

/**
 * AnalysisAverageSalesPerDay
 * @params records []model.Record
 * @return [][]int64
 */
func AnalysisAverageSalesPerDay(records []model.Record) []int64 {
	var mon_sales_index []int64
	var tue_sales_index []int64
	var wed_sales_index []int64
	var thu_sales_index []int64
	var fri_sales_index []int64
	var sat_sales_index []int64
	var sun_sales_index []int64

	var mon_sum int64
	var tue_sum int64
	var wed_sum int64
	var thu_sum int64
	var fri_sum int64
	var sat_sum int64
	var sun_sum int64

	var mon_ave int64
	var tue_ave int64
	var wed_ave int64
	var thu_ave int64
	var fri_ave int64
	var sat_ave int64
	var sun_ave int64

	// 曜日別に振り分け
	for _, record := range records {
		switch record.DayOfWeek {
		case "Mon.":
			mon_sales_index = append(mon_sales_index, record.DailySales)
		case "Tue.":
			tue_sales_index = append(tue_sales_index, record.DailySales)
		case "Wed.":
			wed_sales_index = append(wed_sales_index, record.DailySales)
		case "Thu.":
			thu_sales_index = append(thu_sales_index, record.DailySales)
		case "Fri.":
			fri_sales_index = append(fri_sales_index, record.DailySales)
		case "Sat.":
			sat_sales_index = append(sat_sales_index, record.DailySales)
		case "Sun.":
			sun_sales_index = append(sun_sales_index, record.DailySales)
		}
	}

	// 曜日別平均値の算出
	for _, sale_index := range mon_sales_index {
		mon_sum += sale_index
	}
	mon_ave = mon_sum / int64(len(mon_sales_index))

	for _, sale_index := range tue_sales_index {
		tue_sum += sale_index
	}
	tue_ave = mon_sum / int64(len(tue_sales_index))

	for _, sale_index := range wed_sales_index {
		wed_sum += sale_index
	}
	wed_ave = wed_sum / int64(len(wed_sales_index))

	for _, sale_index := range thu_sales_index {
		thu_sum += sale_index
	}
	thu_ave = thu_sum / int64(len(thu_sales_index))

	for _, sale_index := range fri_sales_index {
		fri_sum += sale_index
	}
	fri_ave = fri_sum / int64(len(fri_sales_index))

	for _, sale_index := range sat_sales_index {
		sat_sum += sale_index
	}
	sat_ave = sat_sum / int64(len(sat_sales_index))

	for _, sale_index := range sun_sales_index {
		sun_sum += sale_index
	}
	sun_ave = sun_sum / int64(len(sun_sales_index))

	return []int64{mon_ave, tue_ave, wed_ave, thu_ave, fri_ave, sat_ave, sun_ave}
}

/**
 * AnalysisAverageOccupancyRatePerDay
 * @params records []model.Record
 * @return [][]float64
 */
func AnalysisAverageOccupancyRatePerDay(records []model.Record) []float64 {
	var mon_occupancy_rate_index []float64
	var tue_occupancy_rate_index []float64
	var wed_occupancy_rate_index []float64
	var thu_occupancy_rate_index []float64
	var fri_occupancy_rate_index []float64
	var sat_occupancy_rate_index []float64
	var sun_occupancy_rate_index []float64

	var mon_sum float64
	var tue_sum float64
	var wed_sum float64
	var thu_sum float64
	var fri_sum float64
	var sat_sum float64
	var sun_sum float64

	var mon_ave float64
	var tue_ave float64
	var wed_ave float64
	var thu_ave float64
	var fri_ave float64
	var sat_ave float64
	var sun_ave float64

	// 曜日別に振り分け
	for _, record := range records {
		switch record.DayOfWeek {
		case "Mon.":
			mon_occupancy_rate_index = append(mon_occupancy_rate_index, record.OccupancyRate)
		case "Tue.":
			tue_occupancy_rate_index = append(tue_occupancy_rate_index, record.OccupancyRate)
		case "Wed.":
			wed_occupancy_rate_index = append(wed_occupancy_rate_index, record.OccupancyRate)
		case "Thu.":
			thu_occupancy_rate_index = append(thu_occupancy_rate_index, record.OccupancyRate)
		case "Fri.":
			fri_occupancy_rate_index = append(fri_occupancy_rate_index, record.OccupancyRate)
		case "Sat.":
			sat_occupancy_rate_index = append(sat_occupancy_rate_index, record.OccupancyRate)
		case "Sun.":
			sun_occupancy_rate_index = append(sun_occupancy_rate_index, record.OccupancyRate)
		}
	}

	// 曜日別平均値の算出
	for _, occupancy_rate_index := range mon_occupancy_rate_index {
		mon_sum += occupancy_rate_index
	}
	mon_ave = mon_sum / float64(len(mon_occupancy_rate_index))

	for _, occupancy_rate_index := range tue_occupancy_rate_index {
		tue_sum += occupancy_rate_index
	}
	tue_ave = mon_sum / float64(len(tue_occupancy_rate_index))

	for _, occupancy_rate_index := range wed_occupancy_rate_index {
		wed_sum += occupancy_rate_index
	}
	wed_ave = wed_sum / float64(len(wed_occupancy_rate_index))

	for _, occupancy_rate_index := range thu_occupancy_rate_index {
		thu_sum += occupancy_rate_index
	}
	thu_ave = thu_sum / float64(len(thu_occupancy_rate_index))

	for _, occupancy_rate_index := range fri_occupancy_rate_index {
		fri_sum += occupancy_rate_index
	}
	fri_ave = fri_sum / float64(len(fri_occupancy_rate_index))

	for _, occupancy_rate_index := range sat_occupancy_rate_index {
		sat_sum += occupancy_rate_index
	}
	sat_ave = sat_sum / float64(len(sat_occupancy_rate_index))

	for _, occupancy_rate_index := range sun_occupancy_rate_index {
		sun_sum += occupancy_rate_index
	}
	sun_ave = sun_sum / float64(len(sun_occupancy_rate_index))

	return []float64{mon_ave, tue_ave, wed_ave, thu_ave, fri_ave, sat_ave, sun_ave}
}

/**
 * PeriodAnalysisData
 * records []model.Record
 * [][]float64
 */
func PeriodAnalysisData(records []model.Record) interface{} {
	var dailyAverageSales int64
	var dailyAverageOccupancyRate float64
	var periodTimeUnitSales int64     // 時間単価
	var periodCustomerUnitSales int64 // 客単価
	var periodDistanceUnitSales int64 // 距離単価
	var maxSales int64                // 最高売上
	var maxOccupancyRate float64      // 最高実車率

	var sales_sum int64
	var occupancy_rate_sum float64
	var running_time_sum int64
	var running_km_sum int64
	var number_of_time int64

	records_number := len(records)

	// 各項目合計値の取得
	for _, record := range records {
		sales_sum += record.DailySales
		occupancy_rate_sum += record.OccupancyRate
		running_time_sum += record.RunningTime
		running_km_sum += record.RunningKm
		number_of_time += record.NumberOfTime

		sales := record.DailySales
		if maxSales <= sales {
			maxSales = sales
		}

		occupancyRate := record.OccupancyRate
		if maxOccupancyRate <= occupancyRate {
			maxOccupancyRate = occupancyRate
		}
	}

	// 平均化
	dailyAverageSales = sales_sum / int64(records_number)
	shift := math.Pow(10, 1)
	dailyAverageOccupancyRate = roundInt(occupancy_rate_sum/float64(records_number)*shift) / shift
	periodTimeUnitSales = sales_sum / running_time_sum
	periodCustomerUnitSales = sales_sum / number_of_time
	periodDistanceUnitSales = sales_sum / running_km_sum

	data := map[string]interface{}{
		"dailyAverageSales":         dailyAverageSales,
		"dailyAverageOccupancyRate": dailyAverageOccupancyRate,
		"periodTimeUnitSales":       periodTimeUnitSales,
		"periodCustomerUnitSales":   periodCustomerUnitSales,
		"periodDistanceUnitSales":   periodDistanceUnitSales,
		"maxSales":                  maxSales,
		"maxOccupancyRate":          math.Round(maxOccupancyRate * 10 / 10),
	}

	return data
}

// 四捨五入用関数
func roundInt(num float64) float64 {
	t := math.Trunc(num)
	if math.Abs(num-t) >= 0.5 {
		return t + math.Copysign(1, num)
	}
	return t
}

package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * Analysis(Home)
 */
func Analysis(c *fiber.Ctx) error {
	log.Println("start analysis (Home)")
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "user not signin")
	}

	// データ収集
	// AnalysisSalesSum
	sales_data_sum, dates_sum, err := service.DataSettingForSalesSum(c)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// AnalysisSales
	sales_data, dates, err := service.GetSalesIndex(c)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// GetRecords
	// データ収集
	records, err := service.SearchRecordForMonth(c)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	data := map[string]interface{}{
		"home_sales_sum": sales_data_sum,
		"dates_sum":      dates_sum,
		"home_sales":     sales_data,
		"dates":          dates,
		"records":        records,
	}

	return service.SuccessResponse(c, []string{"success_get_analysis_data"}, data, nil)
}

/**
 * AnalysisPage
 */
func AnalysisPage(c *fiber.Ctx) error {
	log.Println("start analysis (Analysis)")
	// 曜日別平均の解析
	// 曜日別平均売上
	// 曜日別平均実車率
	// 曜日別平均乗車回数
	// 曜日別平均走行距離
	// 乗車方式別平均売上
	average_sales_per_day, average_occupancy_rate_per_day, err := service.GetAllAnalysisData(c)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return c.JSON(fiber.Map{
		"average_sales_per_day":          average_sales_per_day,
		"average_occupancy_rate_per_day": average_occupancy_rate_per_day,
	})
}

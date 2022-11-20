package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * AnalysisSalesSum
 */
func AnalysisSalesSum(c *fiber.Ctx) error {
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
	// admin権限の確認
	if !statuses[1] {
		// user合致確認
		if !statuses[2] {
			return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
		}
	}

	// データ収集
	service.DataSettingForSalesSum()

	// 仮データ
	labels := []string{"a", "b", "c"}
	data := []int{1, 2, 3}

	return service.SuccessResponseAnalysis(c, []string{"analysis_success"}, data, labels)
}

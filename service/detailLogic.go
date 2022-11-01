package service

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * SearchDetail
 * detailsの検索
 * @params c *fiber.Ctx
 * @returns details []*models.Detail
 */
func SearchDetail(c *fiber.Ctx, adminStatus bool) ([]*model.Detail, error) {
	var details []*model.Detail

	// paramsの確認
	record_id := c.Query("record_id")
	depart_hour := c.Query("depart_hour")
	depart_place := c.Query("depart_place")
	arrive_place := c.Query("arrive_place")
	sales := c.Query("sales")
	method_flg := c.Query("method_flg")
	description := c.Query("description")

	// クエリの作成
	detailSearch := db.DB.Where("")
	if len(depart_hour) != 0 {
		detailSearch.Where("depart_hour = ?", depart_hour)
	}
	if len(depart_place) != 0 {
		detailSearch.Where("depart_place LIKE ?", "%"+depart_place+"%")
	}
	if len(arrive_place) != 0 {
		detailSearch.Where("arrive_place LIKE ?", "%"+arrive_place+"%")
	}
	if len(sales) != 0 {
		detailSearch.Where("sales >= ?", sales)
	}
	if len(method_flg) != 0 {
		detailSearch.Where("method_flg = ?", method_flg)
	}
	if len(description) != 0 {
		detailSearch.Where("description LIKE ?", "%"+description+"%")
	}
	// 対象record_idのdetailsの取得
	detailSearch.Where("record_id = ?", record_id)

	// detailsレコードの取得
	err := detailSearch.Debug().Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

/**
 * DetailValidation
 * detailのバリデーション機能
 * @params detail *model.Detail
 * @returns bool
 * @returns []string
 */
func DetailValidation(detail *model.Detail) (bool, []string) {
	var errs []string

	// depart_hour
	if detail.DepartHour < 0 || 24 < detail.DepartHour {
		log.Println("depart_hour range error")
		errs = append(errs, "depart_hour_range_error")
	}
	// depart_place
	if len(detail.DepartPlace) == 0 {
		log.Println("depart_place null error")
		errs = append(errs, "depart_place_null_error")
	}
	// arrive_place
	if len(detail.DepartPlace) == 0 {
		log.Println("arrive_place null error")
		errs = append(errs, "arrive_place_null_error")
	}
	// sales
	if detail.Sales <= 0 {
		log.Println("sales range error")
		errs = append(errs, "sales_range_error")
	}
	// method_flg
	if len(detail.MethodFlg) == 0 {
		log.Println("method_flg null error")
		errs = append(errs, "method_flg_null_error")
	} else {
		if !(detail.MethodFlg == "flow" || detail.MethodFlg == "wait" || detail.MethodFlg == "app" || detail.MethodFlg == "wireless" || detail.MethodFlg == "own" || detail.MethodFlg == "other") {
			log.Println("specified word error(method_flg)")
			errs = append(errs, "specified_word_error(method_flg)")
		}
	}

	// errの出力
	if len(errs) != 0 {
		return false, errs
	}

	return true, nil
}

/**
 * GetDetail
 * detail情報の取得
 * @params c *fiber.Ctx
 * @returns detail *model.Detail
 */
func GetDetail(c *fiber.Ctx) (*model.Detail, error) {
	// 変数確認
	detail_id := c.Params("id")
	var detail *model.Detail

	// レコードの取得
	err := db.DB.Where("id = ?", detail_id).First(&detail).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return nil, fmt.Errorf(constants.DB_ERR)
	}

	return detail, nil
}

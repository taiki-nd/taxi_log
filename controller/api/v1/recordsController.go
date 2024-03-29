package controller

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * RecordsIndex
 * recordの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsIndex(c *fiber.Ctx) error {
	log.Println("start RecordsIndex")
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

	// recordの検索
	records, err := service.SearchRecord(c, statuses[1])
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// pagination
	var total int64
	db.DB.Model(&model.Record{}).Count(&total)
	limit := constants.PAGINATION_LIMIT
	page, _ := strconv.Atoi(c.Query("page", "1"))
	lastPage := math.Ceil(float64(total) / float64(limit))
	meta := &model.Meta{
		Total:    total,
		Page:     page,
		LastPage: lastPage,
	}

	return service.SuccessResponse(c, []string{"index_record_success"}, records, meta)
}

/**
 * RecordsShow
 * recordの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsShow(c *fiber.Ctx) error {
	log.Println("start RecordsShow")
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
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	// レコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// admin権限の確認
	/*
		if !statuses[1] {
			// follower確認
			status, err := service.IsFollowerForRecord(c, record)
			if err != nil {
				return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
			}
			if !status {
				return service.ErrorResponse(c, []string{constants.FOLLOW_RELATIONSHIP_ERROR}, "follow relationship error")
			}
		}
	*/

	return service.SuccessResponse(c, []string{"show_record_success"}, record, nil)
}

/**
 * RecordsCreate
 * recordの新規登録処理
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsCreate(c *fiber.Ctx) error {
	log.Println("start RecordsCreate")
	// 変数確認
	var record *model.Record

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

	// リクエストボディーのパース
	err = c.BodyParser(&record)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, []string{constants.BODY_PARSE_ERROR}, fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs = service.RecordValidation(record)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// レコード作成
	err = db.DB.Create(&record).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"create_record_success"}, record, nil)
}

/**
 * RecordsUpdate
 * record情報の更新処理
 * @params c *fiber.Ctx
 */
func RecordsUpdate(c *fiber.Ctx) error {
	log.Println("start RecordsUpdate")
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
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "record not signin")
	}
	// user合致確認
	if !statuses[2] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_MATCH}, "user not match")
	}

	// recordレコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// リクエストボディのパース
	err = c.BodyParser(record)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return service.ErrorResponse(c, []string{constants.BODY_PARSE_ERROR}, fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs = service.RecordValidation(record)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	update_record := map[string]interface{}{
		"id":                   record.Id,
		"date":                 record.Date,
		"day_of_week":          record.DayOfWeek,
		"style_flg":            record.StyleFlg,
		"start_hour":           record.StartHour,
		"running_time":         record.RunningTime,
		"running_km":           record.RunningKm,
		"occupancy_rate":       record.OccupancyRate,
		"number_of_time":       record.NumberOfTime,
		"is_tax":               record.IsTax,
		"daily_sales":          record.DailySales,
		"daily_sales_with_tax": record.DailySalesWithTax,
		"user_id":              record.UserId,
	}

	// record情報の更新
	err = db.DB.Model(&record).Updates(update_record).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	return service.SuccessResponse(c, []string{"update_record_success"}, record, nil)
}

/**
 * RecordsDelete
 * record情報の削除
 * @params c *fiber.Ctx
 * @returns error
 */
func RecordsDelete(c *fiber.Ctx) error {
	log.Println("start RecordsDelete")
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

	// recordレコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}

	// record情報の削除トランザクション開始
	tx := db.DB.Begin()
	err = service.RecordDeleteTransaction(tx, record)
	if err != nil {
		log.Printf("db error: %v", err)
		return service.ErrorResponse(c, []string{constants.DB_ERR}, fmt.Sprintf("db error: %v", err))
	}
	tx.Commit()

	return service.SuccessResponse(c, []string{"create_delete_success"}, nil, nil)
}

package controller

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/service"
)

/**
 * RecordsIndex
 * recordの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsIndex(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_auth_error",
				"message": fmt.Sprintf("user auth error: %v", err),
			},
			"data": fiber.Map{},
		})
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_not_signin",
				"message": "user not signin",
			},
			"data": fiber.Map{},
		})
	}

	// recordの検索
	records, err := service.SearchRecord(c, statuses[1])
	if err != nil {
		log.Printf("db error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "db_error",
				"message": fmt.Sprintf("db error: %v", err),
			},
			"data": fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  true,
			"code":    "index_record_success",
			"message": "",
		},
		"data": records,
	})
}

/**
 * RecordsShow
 * recordの一覧取得
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsShow(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_auth_error",
				"message": fmt.Sprintf("user auth error: %v", err),
			},
			"data": fiber.Map{},
		})
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_not_signin",
				"message": "user not signin",
			},
			"data": fiber.Map{},
		})
	}
	// user合致確認
	if !statuses[2] {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_not_match",
				"message": "user not match",
			},
			"data": fiber.Map{},
		})
	}

	// レコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "db_error",
				"message": fmt.Sprintf("db error: %v", err),
			},
			"data": fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  true,
			"code":    "show_record_success",
			"message": "",
		},
		"data": record,
	})
}

/**
 * RecordsCreate
 * recordの新規登録処理
 * @params c *fiber.Ctx
 * @returns error error
 */
func RecordsCreate(c *fiber.Ctx) error {
	// 変数確認
	var record *model.Record

	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_auth_error",
				"message": fmt.Sprintf("user auth error: %v", err),
			},
			"data": fiber.Map{},
		})
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_not_signin",
				"message": "user not signin",
			},
			"data": fiber.Map{},
		})
	}

	// リクエストボディーのパース
	err = c.BodyParser(&record)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "body_parse_error",
				"message": fmt.Sprintf("body parse error: %v", err),
			},
			"data": fiber.Map{},
		})
	}

	// バリデーション
	_, errs = service.RecordValidation(record)
	if len(errs) != 0 {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    errs,
				"message": fmt.Sprintf("validation error: %v", errs),
			},
			"data": fiber.Map{},
		})
	}

	// レコード作成
	err = db.DB.Create(&record).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "db_error",
				"message": fmt.Sprintf("db error: %v", err),
			},
			"data": fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"info": fiber.Map{
			"status":  true,
			"code":    "create_record_success",
			"message": "",
		},
		"data": record,
	})
}

/**
 * RecordsUpdate
 * record情報の更新処理
 * @params c *fiber.Ctx
 */
func RecordsUpdate(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_auth_error",
				"message": fmt.Sprintf("user auth error: %v", err),
			},
			"data": fiber.Map{},
		})
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "record_not_signin",
				"message": "record not signin",
			},
			"data": fiber.Map{},
		})
	}
	// user合致確認
	if !statuses[2] {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_not_match",
				"message": "user not match",
			},
			"data": fiber.Map{},
		})
	}

	// recordレコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return c.JSON(fiber.Map{
			"status": false,
			"code":   "db_error",
			"data":   fiber.Map{},
		})
	}

	// リクエストボディのパース
	err = c.BodyParser(record)
	if err != nil {
		log.Printf("body parse error: %v", err)
		return c.JSON(fiber.Map{
			"status": false,
			"code":   "body_parse_error",
			"data":   fiber.Map{},
		})
	}

	// バリデーション
	_, errs = service.RecordValidation(record)
	if len(errs) != 0 {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    errs,
				"message": fmt.Sprintf("validation error: %v", errs),
			},
			"data": fiber.Map{},
		})
	}

	// record情報の更新
	err = db.DB.Model(&record).Updates(record).Error
	if err != nil {
		log.Printf("db error: %v", err)
		return c.JSON(fiber.Map{
			"status": false,
			"code":   "db_error",
			"data":   fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"status": true,
		"code":   "update_record_success",
		"data":   record,
	})
}

/**
 * RecordsDelete
 * record情報の削除
 * @params c *fiber.Ctx
 * @returns error
 */
func RecordsDelete(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_auth_error",
				"message": fmt.Sprintf("user auth error: %v", err),
			},
			"data": fiber.Map{},
		})
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return c.JSON(fiber.Map{
			"info": fiber.Map{
				"status":  false,
				"code":    "user_not_signin",
				"message": "user not signin",
			},
			"data": fiber.Map{},
		})
	}
	// admin権限の確認
	if !statuses[1] {
		// user合致確認
		if !statuses[2] {
			return c.JSON(fiber.Map{
				"info": fiber.Map{
					"status":  false,
					"code":    "user_not_match",
					"message": "user not match",
				},
				"data": fiber.Map{},
			})
		}
	}

	// recordレコードの取得
	record, err := service.GetRecord(c)
	if err != nil {
		log.Printf("db error: %v", err)
		return c.JSON(fiber.Map{
			"status": false,
			"code":   "db_error",
			"data":   fiber.Map{},
		})
	}

	// record情報の削除
	errRecord := db.DB.Delete(record).Error
	if errRecord != nil {
		log.Printf("db error: %v", err)
		return c.JSON(fiber.Map{
			"status": false,
			"code":   "db_error",
			"data":   fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"status": true,
		"code":   "delete_record_success",
		"data":   fiber.Map{},
	})
}

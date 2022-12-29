package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
)

/**
 * GetAreasFromPrefecture
 */
func GetAreasFromPrefecture(c *fiber.Ctx) ([]string, error) {
	prefecture := c.Query("prefecture")

	// area一覧の取得
	var areas []string
	err := db.DB.Table("areas").Where("prefecture = ?", prefecture).Pluck("area", &areas).Error
	if err != nil {
		return nil, err
	}
	return areas, nil
}

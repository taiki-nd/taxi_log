package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
)

type Follow struct {
	Id              uint
	Nickname        string
	ProfileImageUrl string
}

/**
 * IsFollower
 * follower確認
 * @params c *fiber.Ctx
 * @returns bool
 */
func IsFollower(c *fiber.Ctx) (bool, error) {
	// 変数確認
	signin_user, _ := GetUserFromUuid(c)
	id := c.Params("id")
	var follow model.UserFollowing

	// follow関係の確認
	err := db.DB.Table("user_followings").Where("user_id = ?", signin_user.Id).Where("following_id = ?", id).First(&follow).Error
	if err != nil {
		return false, fmt.Errorf("follow_relationship_error: %v", err)
	} else {
		return true, nil
	}
}

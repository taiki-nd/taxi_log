package service

import (
	"fmt"
	"log"

	"github.com/taiki-nd/taxi_log/model"
	"gorm.io/gorm"
)

func UsersDeleteTransaction(tx *gorm.DB, user *model.User) error {
	// recordの削除
	errRecord := tx.Table("records").Where("user_id = ?", user.Id).Delete("").Error
	if errRecord != nil {
		log.Printf("transaction delete record err: %v", errRecord)
		return fmt.Errorf("transaction_delete_record_err")
	}
	// follow関連の削除
	errFollow := tx.Table("user_followings").Where("user_id = ?", user.Id).Or("following_id = ?", user.Id).Delete("").Error
	if errFollow != nil {
		log.Printf("transaction delete follow err: %v", errFollow)
		return fmt.Errorf("transaction_delete_follow_err")
	}
	// userの削除
	errUser := tx.Delete(user).Error
	if errUser != nil {
		log.Printf("transaction delete user err: %v", errUser)
		return fmt.Errorf("transaction_delete_user_err")
	}
	return nil
}

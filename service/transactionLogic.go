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
	// userの削除
	errUser := tx.Delete(user).Error
	if errUser != nil {
		log.Printf("transaction delete user err: %v", errRecord)
		return fmt.Errorf("transaction_delete_user_err")
	}
	return nil
}

package db

import (
	"fmt"
	"log"

	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

/**
 * ConnectToDb
 * db接続の設定
 */
func ConnectToDb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		config.Config.User, config.Config.Password,
		config.Config.Host, config.Config.Port,
		config.Config.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	log.Println("db connection success")

	DB = db

	DB.AutoMigrate(
		model.User{},
		model.Detail{},
		model.Record{},
		model.UserFollowing{},
	)

}

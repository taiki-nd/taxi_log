package main

import (
	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/utils"
)

func main() {
	// logの有効化
	utils.Logging(config.Config.LogFile)

	// db接続
	db.ConnectToDb()
}

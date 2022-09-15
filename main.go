package main

import (
	"log"

	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/utils"
)

func main() {
	// logの有効化
	utils.Logging(config.Config.LogFile)
	log.Println("Hello World!")
}

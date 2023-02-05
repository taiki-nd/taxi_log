package crons

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func CronManager(cron *cron.Cron) {
	cron.AddFunc("@every 3s", Hello)
}

func Hello() {
	fmt.Println("hello")
}

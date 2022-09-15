package utils

import (
	"io"
	"log"
	"os"
	"time"
)

/*
 * logging
 * 出力ログファイル・出力ログの設定
 * @params logFile string
 */
func Logging(logFile string) {
	// 日付とフォーマットの設定
	day := time.Now()
	const layout = "2006-01-02"

	nowUTC := day.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	dayJST := nowUTC.In(jst)

	// logファイルの生成
	log_file, err := os.OpenFile(logFile+"_"+dayJST.Format(layout)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=logFile err=%s", err.Error())
	}
	multiLogFile := io.MultiWriter(os.Stdout, log_file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)
}

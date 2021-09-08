package main

import (
	"os"
	"time"
	"trip-weather-backend/config"
	"trip-weather-backend/model"
	"trip-weather-backend/route"

	log "github.com/sirupsen/logrus"
)

func init() {
	// ロギングを初期設定
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	// ログ出力設定(標準出力)
	log.SetOutput((os.Stdout))
	// ログ出力設定(ログファイル出力)
	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 	log.SetOutput(file)
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }
}

func main() {
	// DBを初期化
	model.InitDB(config.DB_PATH)

	// for Heroku 環境変数からポート番号を取得
	var port string = os.Getenv("PORT")
	// for Local 環境変数に"PORT"がない場合、configの値を利用
	if port == "" {
		port = config.PORT
	}

	// routerの初期化、起動
	router := route.InitRoute()
	router.Logger.Fatal(router.Start(":" + port))
}

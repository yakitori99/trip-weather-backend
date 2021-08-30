package main

import (
	"os"
	"trip-weather-backend/config"
	"trip-weather-backend/model"
	"trip-weather-backend/route"
)

func main() {
	// DBを初期化
	model.InitDB()

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

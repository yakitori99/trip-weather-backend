package main

import (
	"trip-weather-backend/config"
	"trip-weather-backend/model"
	"trip-weather-backend/route"
)

func main() {
	// DBを初期化
	model.InitDB()

	// routerの初期化、起動
	router := route.InitRoute()
	router.Logger.Fatal(router.Start(":" + config.PORT))
}

package model

import (
	// "gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"trip-weather-backend/utils"
)

// パッケージ内で利用するdb(gorm.DBのインスタンスのポインタ)を宣言
var db *gorm.DB

// DBファイルのパスを受け取り、DBを初期化する
func InitDB(dsn string) {
	// データベースをインスタンス化し、パッケージ内変数のdb(変数)に代入
	// 一旦一時的な変数に入れ、2行目でパッケージ内変数dbへ代入する（こうしないとパッケージ内変数へ代入されない）
	// _db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db = _db

	if err != nil {
		utils.OutErrorLog("failed to connect database", err)
		panic("failed to connect database")
	}

	// Migrate the schema
	// db.AutoMigrate(&Pref{})
	// db.AutoMigrate(&City{})
	// db.AutoMigrate(&Favorite{})

	utils.OutInfoLog("END")
}

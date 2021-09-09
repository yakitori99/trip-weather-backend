package model

import (
	"time"
	"trip-weather-backend/config"

	"gorm.io/gorm"
)

type Favorite struct {
	// お気に入りテーブル
	// 論理削除でなく、物理削除を使いたいのでgorm.Modelは利用しない
	ID           int `gorm:"primary_key"`
	FromPrefCode string
	FromCityCode string
	ToPrefCode   string
	ToCityCode   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type Favorites []Favorite

// favoritesテーブルに対し、cityCodeで見て同一レコードがあればUPD, なければINSする関数
func CreateFavoriteTransaction(fromCityCode string, toCityCode string) (int, error) {
	var resultCode int
	// トランザクション開始
	tx := db.Begin()

	// SELECTして同一レコードが存在するかチェック
	var favorite Favorite
	err := tx.Where("from_city_code = ? AND to_city_code = ?", fromCityCode, toCityCode).First(&favorite).Error
	if err == gorm.ErrRecordNotFound {
		// 同一レコードがなければINSERT
		err = CreateFavorite(tx, fromCityCode, toCityCode)
		resultCode = config.DONE_INS
	} else if err == nil {
		// 同一レコードがあればUpdatedAtのみUpdate
		err = tx.Model(Favorite{}).Where("from_city_code = ? AND to_city_code = ?", fromCityCode, toCityCode).Updates(Favorite{UpdatedAt: time.Now()}).Error
		resultCode = config.DONE_UPD
	}

	if err != nil {
		tx.Rollback()
		return config.DONE_ERR, err
	}

	tx.Commit()
	return resultCode, nil
}

// favoriteテーブルにレコードをINSERTする関数
func CreateFavorite(tx *gorm.DB, fromCityCode string, toCityCode string) error {
	favorite := Favorite{
		FromPrefCode: fromCityCode[:2],
		FromCityCode: fromCityCode,
		ToPrefCode:   toCityCode[:2],
		ToCityCode:   toCityCode,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	result := tx.Create(&favorite)
	return result.Error
}

// favoriteテーブルからレコードをDELETEする関数
func DeleteFavorite(fromCityCode string, toCityCode string) error {
	result := db.Where("from_city_code = ? AND to_city_code = ?", fromCityCode, toCityCode).Delete(Favorite{})
	return result.Error
}

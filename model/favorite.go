package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"trip-weather-backend/config"
	"trip-weather-backend/utils"

	"gorm.io/gorm"
)

type Favorite struct {
	// お気に入りテーブル
	// 論理削除でなく、物理削除を使いたいのでgorm.Modelは利用しない
	ID           int `gorm:"primary_key"`
	Nickname     string
	FromPrefCode string
	FromCityCode string
	ToPrefCode   string
	ToCityCode   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type Favorites []Favorite

// 画面表示に必要な要素を追加・抜粋したFavorite
type SelectedFavorite struct {
	Nickname     string
	FromPrefCode string
	FromCityCode string
	ToPrefCode   string
	ToCityCode   string
	FromPrefName string
	FromCityName string
	ToPrefName   string
	ToCityName   string
	UpdatedAt    time.Time
}
type SelectedFavorites []SelectedFavorite

// ニックネーム一覧の返却用構造体
type SelectedNickname struct {
	Nickname string
}
type SelectedNicknames []SelectedNickname

var getFavoriteSql string = `select
f.nickname, f.from_pref_code, f.from_city_code, f.to_pref_code,f.to_city_code,
p1.pref_name as from_pref_name, 
c1.city_name as from_city_name, 
p2.pref_name as to_pref_name, 
c2.city_name as to_city_name, 
f.updated_at
from favorites f
LEFT OUTER JOIN prefs p1 on f.from_pref_code = p1.pref_code
LEFT OUTER JOIN prefs p2 on f.to_pref_code = p2.pref_code 
LEFT OUTER JOIN cities c1 on f.from_city_code = c1.city_code
LEFT OUTER JOIN cities c2 on f.to_city_code = c2.city_code 
order by f.updated_at desc`

var getFavoriteByNicknameSql string = `select
f.nickname, f.from_pref_code, f.from_city_code, f.to_pref_code,f.to_city_code,
p1.pref_name as from_pref_name, 
c1.city_name as from_city_name, 
p2.pref_name as to_pref_name, 
c2.city_name as to_city_name, 
f.updated_at
from favorites f
LEFT OUTER JOIN prefs p1 on f.from_pref_code = p1.pref_code
LEFT OUTER JOIN prefs p2 on f.to_pref_code = p2.pref_code 
LEFT OUTER JOIN cities c1 on f.from_city_code = c1.city_code
LEFT OUTER JOIN cities c2 on f.to_city_code = c2.city_code 
where f.nickname = ? 
order by f.updated_at desc`

// favoritesテーブルに対し、nicknameとcityCodeで見て同一レコードがあればUPD, なければINSする関数
func CreateFavoriteTransaction(nickname string, fromCityCode string, toCityCode string) (int, error) {
	var resultCode int
	// 空白文字は削除(半角/全角とも)
	nickname = strings.ReplaceAll(nickname, " ", "")
	nickname = strings.ReplaceAll(nickname, "　", "")

	// トランザクション開始
	tx := db.Begin()

	// fromCityCode , toCityCodeの片方でもcitiesテーブルに存在しない場合はエラー
	var count int64
	for _, cityCode := range []string{fromCityCode, toCityCode} {
		tx.Model(&City{}).Where("city_code = ?", cityCode).Count(&count)
		if count == 0 {
			tx.Rollback()
			err := errors.New("invalid city_code")
			utils.OutErrorLogDetail("invalid city_code", err, fmt.Sprintf("invalid cityCode:%v", cityCode))
			return config.DONE_ERR, err
		}
	}

	// SELECTして同一レコードが存在するかチェック
	var favorite Favorite
	err := tx.Where("nickname = ? AND from_city_code = ? AND to_city_code = ?", nickname, fromCityCode, toCityCode).First(&favorite).Error
	if err == gorm.ErrRecordNotFound {
		// 同一レコードがなければINSERT
		err = CreateFavorite(tx, nickname, fromCityCode, toCityCode)
		resultCode = config.DONE_INS
	} else if err == nil {
		// 同一レコードがあればUpdatedAtのみUpdate
		err = tx.Model(Favorite{}).Where("nickname = ? AND from_city_code = ? AND to_city_code = ?", nickname, fromCityCode, toCityCode).Updates(Favorite{UpdatedAt: time.Now()}).Error
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
func CreateFavorite(tx *gorm.DB, nickname string, fromCityCode string, toCityCode string) error {
	favorite := Favorite{
		Nickname:     nickname,
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
func DeleteFavorite(nickname string, fromCityCode string, toCityCode string) error {
	result := db.Where("nickname = ? AND from_city_code = ? AND to_city_code = ?", nickname, fromCityCode, toCityCode).Delete(Favorite{})
	return result.Error
}

// favoritesテーブルから更新日時降順で全件を取得し、pref_name, city_nameとJOINして返す関数
func GetFavoriteAll() SelectedFavorites {
	var selectedFavorites SelectedFavorites
	db.Raw(getFavoriteSql).Scan(&selectedFavorites)
	return selectedFavorites
}

// favoritesテーブルから更新日時降順でn件を取得し、pref_name, city_nameとJOINして返す関数
func GetFavoriteN(n int) SelectedFavorites {
	var selectedFavorites SelectedFavorites
	// n件まで の構文を追加
	rawSql := getFavoriteSql + " limit ?"
	db.Raw(rawSql, n).Scan(&selectedFavorites)
	return selectedFavorites
}

// favoritesテーブルからNicknameをキーに更新日時降順で全件を取得し、pref_name, city_nameとJOINして返す関数
func GetFavoriteByNickname(nickname string) SelectedFavorites {
	var selectedFavorites SelectedFavorites
	db.Raw(getFavoriteByNicknameSql, nickname).Scan(&selectedFavorites)
	return selectedFavorites
}

// favoritesテーブルから重複しないNickname一覧を取得し、Nicknameの昇順で返す関数
func GetNicknameDistinct() SelectedNicknames {
	var selectedNicknames SelectedNicknames
	db.Table("favorites").Distinct().Order("nickname asc").Pluck("nickname", &selectedNicknames)
	return selectedNicknames
}

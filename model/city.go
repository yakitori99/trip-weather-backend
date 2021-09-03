package model

import (
	"time"
)

type City struct {
	// 都市テーブル
	// 論理削除でなく、物理削除を使いたいのでgorm.Modelは利用しない
	ID                 int `gorm:"primary_key"`
	CityCode           string
	CityName           string
	PrefCode           string
	CityKana           string
	CityRomaji         string
	CityRomajiLocation string
	CityLon            float64
	CityLat            float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
type Cities []City

// 要素を抜粋したCity
type SelectedCity struct {
	CityCode string
	CityName string
	PrefCode string
	CityLon  float64
	CityLat  float64
}
type SelectedCities []SelectedCity

// Citiesテーブルから全行の抜粋した列を取得し、city_codeの昇順でソートし、返す関数
func GetCityAll() SelectedCities {
	var cities Cities
	db.Order("city_code").Find(&cities)
	// citiesから必要な列のみ抽出
	var selectedCities SelectedCities
	for _, v := range cities {
		selectedCity := SelectedCity{v.CityCode, v.CityName, v.PrefCode, v.CityLon, v.CityLat}
		selectedCities = append(selectedCities, selectedCity)
	}
	return selectedCities
}

// Citiesテーブルからpref_codeで行を抽出し、抜粋した列を取得し、city_codeの昇順でソートし、返す関数
func GetCityByPrefCode(pref_code string) SelectedCities {
	var cities Cities
	db.Order("city_code").Where("pref_code = ?", pref_code).Find(&cities)
	// citiesから必要な列のみ抽出
	var selectedCities SelectedCities
	for _, v := range cities {
		selectedCity := SelectedCity{v.CityCode, v.CityName, v.PrefCode, v.CityLon, v.CityLat}
		selectedCities = append(selectedCities, selectedCity)
	}
	return selectedCities
}

// Citiesテーブルからcity_codeの完全一致で1行を取得し、抜粋した列を返す関数
func GetLocationByCityCode(city_code string) SelectedCity {
	var city City
	db.Where("city_code = ?", city_code).First(&city)
	// cityから必要な列のみ抽出
	selectedCity := SelectedCity{city.CityCode, city.CityName, city.PrefCode, city.CityLon, city.CityLat}
	return selectedCity
}

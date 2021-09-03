package model

import (
	"time"
)

type Pref struct {
	// 都道府県テーブル
	// 論理削除でなく、物理削除を使いたいのでgorm.Modelは利用しない
	ID        int `gorm:"primary_key"`
	PrefCode  string
	PrefName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Prefs []Pref

// 要素を抜粋したPref
type SelectedPref struct {
	PrefCode string
	PrefName string
}
type SelectedPrefs []SelectedPref

func GetPrefAll() SelectedPrefs {
	var prefs Prefs
	db.Select("pref_code, pref_name").Order("pref_code").Find(&prefs)
	// prefsから必要な要素のみ抽出
	var selectedPrefs SelectedPrefs
	for _, v := range prefs {
		selectedPrefs = append(selectedPrefs, SelectedPref{v.PrefCode, v.PrefName})
	}
	return selectedPrefs
}

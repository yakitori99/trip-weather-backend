package model_test

import (
	"os"
	"testing"
	"trip-weather-backend/config"
	"trip-weather-backend/model"
)

func TestMain(m *testing.M) {
	//// 共通前処理
	println("Test START")
	// DB初期化
	dbPath := "../db/trip_weather.db" //UnitTest用のDBパス
	model.InitDB(dbPath)

	//// テスト実行
	code := m.Run()

	//// 共通後処理
	println("Test END")

	os.Exit(code)
}

func Test_CreateFavoriteTransaction(t *testing.T) {
	testSlice := [][]string{
		// INS
		{"ryo", "270000", "130010"},
		{"", "011000", "020010"},
		{"ryo", "130010", "016010"},
		{"", "040010", "400010"},
		// UPD
		{"ryo", "270000", "130010"},
		{"", "011000", "020010"},
		{" ryo　", "270000", "130010"},
		{" r y o　", "270000", "130010"},
		{" r　yo　", "270000", "130010"},
		{" ", "011000", "020010"},
		{"　", "011000", "020010"},
		// ERR
		{"u1", "990020", "130010"},
		{"", "990020", "130010"},
	}

	// テスト準備 レコードがすでに存在する場合DELETE
	for i := range testSlice {
		err := model.DeleteFavorite(testSlice[i][0], testSlice[i][1], testSlice[i][2])
		if err != nil {
			t.Error(err)
		}
	}

	// テスト
	for i := range testSlice {
		resultCode, err := model.CreateFavoriteTransaction(testSlice[i][0], testSlice[i][1], testSlice[i][2])
		if i <= 5 && err != nil {
			t.Error(err)
		}
		// 0-3はINS, 4-10はUPD, 11-はエラーならOK
		if i >= 0 && i <= 3 {
			if resultCode == config.DONE_INS {
				t.Logf("Record:%v, Inserted. OK", i)
			} else {
				t.Errorf("Record:%v, NG", i)
			}
		} else if i >= 4 && i <= 10 {
			if resultCode == config.DONE_UPD {
				t.Logf("Record:%v, Updated. OK", i)
			} else {
				t.Errorf("Record:%v, NG", i)
			}
		} else if i >= 11 {
			if err != nil {
				t.Logf("Record:%v, ERR. OK. err is%v", i, err)
			}
		}

	}

}

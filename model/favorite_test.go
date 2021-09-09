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
	cityCodesSlice := [][]string{
		{"011000", "130010"},
		{"130010", "130010"},
		{"130020", "130010"},
	}

	// テスト準備 1,2は、存在する場合DELETE
	for i := 1; i < 3; i++ {
		err := model.DeleteFavorite(cityCodesSlice[i][0], cityCodesSlice[i][1])
		if err != nil {
			t.Error(err)
		}
	}

	// テスト
	for i := range cityCodesSlice {
		resultCode, err := model.CreateFavoriteTransaction(cityCodesSlice[i][0], cityCodesSlice[i][1])
		if err != nil {
			t.Error(err)
		}
		// 0はUPD, 1-2はINSならOK
		if i == 0 {
			if resultCode == config.DONE_UPD {
				t.Logf("Record:%v, Updated. OK", i)
			} else {
				t.Errorf("Record:%v, NG", i)
			}
		} else if i >= 1 {
			if resultCode == config.DONE_INS {
				t.Logf("Record:%v, Inserted. OK", i)
			} else {
				t.Errorf("Record:%v, NG", i)
			}
		}

	}

}

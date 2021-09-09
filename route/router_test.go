package route_test

import (
	"net/http"
	"os"
	"testing"

	// apitest
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	// 自作パッケージをインポート

	"trip-weather-backend/config"
	"trip-weather-backend/model"
	"trip-weather-backend/route"
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

func Test_CreateFavoriteFromJson(t *testing.T) {
	var favoriteJSON = `{"from_city_code":"170010","to_city_code":"180010"}`
	var favoriteJSONNG1 = `{"fromCityCode":"170010","toCityCode":"180010"}`
	var favoriteJSONNG2 = `{"from_city_code":"170010","to_city_code":"990010"}`

	// テーブル準備
	err := model.DeleteFavorite("170010", "180010")
	if err != nil {
		t.Error(err)
	}

	e := route.InitRoute()

	// 正常系 INS
	apitest.New().
		Handler(e).
		Post("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		JSON(favoriteJSON).
		Expect(t).
		Status(http.StatusCreated).
		Assert(jsonpath.Equal(`$.ResultCode`, float64(config.DONE_INS))).
		End()
	t.Log("INS END")

	// 正常系 UPD
	apitest.New().
		Handler(e).
		Post("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		JSON(favoriteJSON).
		Expect(t).
		Status(http.StatusCreated).
		Assert(jsonpath.Equal(`$.ResultCode`, float64(config.DONE_UPD))).
		End()
	t.Log("UPD END")

	// 異常系1
	apitest.New().
		Handler(e).
		Post("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		JSON(favoriteJSONNG1).
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
	t.Log("ERR1 END")

	// 異常系2
	apitest.New().
		Handler(e).
		Post("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		JSON(favoriteJSONNG2).
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
	t.Log("ERR2 END")
}

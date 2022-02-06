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
	model.InitDB(config.DSN_FOR_TEST)

	//// テスト実行
	code := m.Run()

	//// 共通後処理
	println("Test END")

	os.Exit(code)
}

func Test_CreateFavoriteFromJson(t *testing.T) {
	var favoriteJSON1 = `{"nickname":"",   "from_city_code":"170010","to_city_code":"180010"}`
	var favoriteJSON2 = `{"nickname":"たなか", "from_city_code":"130010","to_city_code":"130010"}`
	var favoriteJSONNG1 = `{"nickname":"", "fromCityCode":"170010",  "toCityCode":"180010"}`
	var favoriteJSONNG2 = `{"nickname":"", "from_city_code":"170010","to_city_code":"990010"}`

	// テーブル準備
	deleteSlice := [][]string{
		{"", "170010", "180010"},
		{"u1", "170010", "180010"},
		{"たなか", "130010", "130010"},
	}
	for _, v := range deleteSlice {
		err := model.DeleteFavorite(v[0], v[1], v[2])
		if err != nil {
			t.Error(err)
		}
	}

	e := route.InitRoute()

	// 正常系 INS
	apitest.New().
		Handler(e).
		Post("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		JSON(favoriteJSON1).
		Expect(t).
		Status(http.StatusCreated).
		Assert(jsonpath.Equal(`$.ResultCode`, float64(config.DONE_INS))).
		End()
	t.Log("INS END")

	apitest.New().
		Handler(e).
		Post("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		JSON(favoriteJSON2).
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
		JSON(favoriteJSON1).
		Expect(t).
		Status(http.StatusCreated).
		Assert(jsonpath.Equal(`$.ResultCode`, float64(config.DONE_UPD))).
		End()
	t.Log("UPD END")

	apitest.New().
		Handler(e).
		Post("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		JSON(favoriteJSON2).
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

func Test_GetFavoriteAll(t *testing.T) {
	e := route.InitRoute()

	// 正常系
	apitest.New().
		Handler(e).
		Get("/favorites").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusOK).
		End()
	t.Log("GET END")
}

func Test_GetFavoriteN(t *testing.T) {
	e := route.InitRoute()

	// 正常系
	apitest.New().
		Handler(e).
		Get("/favorites/to/10").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusOK).
		End()
	t.Log("GET END")

	// 異常系1
	apitest.New().
		Handler(e).
		Get("/favorites/to/a").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
	t.Log("ERR1 END")

	// 異常系2
	apitest.New().
		Handler(e).
		Get("/favorites/to/0").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
	t.Log("ERR2 END")

	// 異常系3
	apitest.New().
		Handler(e).
		Get("/favorites/to/-1").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
	t.Log("ERR3 END")
}

func Test_GetFavoriteByNickname(t *testing.T) {
	e := route.InitRoute()

	// 正常系1
	apitest.New().
		Handler(e).
		Get("/favorites/by/u1").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusOK).
		End()
	t.Log("GET1 END")

	// 正常系2
	apitest.New().
		Handler(e).
		Get("/favorites/by/1").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusOK).
		End()
	t.Log("GET2 END")

	//  正常系3
	apitest.New().
		Handler(e).
		Get("/favorites/by/no_nickname_selected").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusOK).
		End()
	t.Log("GET3 END")

	// 異常系1
	apitest.New().
		Handler(e).
		Get("/favorites/by/").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusNotFound).
		End()
	t.Log("ERR1 END")
}

func Test_GetNicknameDistinct(t *testing.T) {
	e := route.InitRoute()

	// 正常系
	apitest.New().
		Handler(e).
		Get("/nicknames").
		Headers(map[string]string{"User-Agent": "apitest"}).
		Expect(t).
		Status(http.StatusOK).
		End()
	t.Log("GET END")
}

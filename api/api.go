package api

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	echo "github.com/labstack/echo/v4"
	// 自作パッケージをインポート
	"trip-weather-backend/config"
	"trip-weather-backend/model"
	"trip-weather-backend/utils"
)

// FavoriteテーブルへのINSに必須の要素
type FavoriteRequired struct {
	// 受け取ったjsonからの紐付け用に、jsonタグも付ける
	Nickname     string `json:"nickname"`
	FromCityCode string `json:"from_city_code"`
	ToCityCode   string `json:"to_city_code"`
}

// FavoriteテーブルへのINS結果(レスポンス用)
type FavoriteInsResult struct {
	ResultCode int
}

// 画面表示に必要な型変換済みのSelectedFavorite
type SelectedFavoriteStr struct {
	Nickname     string
	FromPrefCode string
	FromCityCode string
	ToPrefCode   string
	ToCityCode   string
	FromPrefName string
	FromCityName string
	ToPrefName   string
	ToCityName   string
	UpdatedAt    string
}
type SelectedFavoriteStrs []SelectedFavoriteStr

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		return c.String(http.StatusOK, "Hello World!")
	}
}

func HelloUsername() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		username := c.Param("username")
		return c.String(http.StatusOK, "Hello World! Your name is "+username)
	}
}

// today-1 ~ +7のdatetimeをstrにしたデータを返す
func GetDatetimes() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		timeToday := time.Now()
		// today-1 ~ +7
		var timeStrDays []string
		for i := -1; i <= 7; i++ {
			timeDay := timeToday.AddDate(0, 0, i)
			timeStrDays = append(timeStrDays, timeDay.Format(config.WEATHER_DATE_FORMAT))
		}
		return c.JSON(http.StatusOK, timeStrDays)
	}
}

func GetPrefs() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		prefs := model.GetPrefAll()
		return c.JSON(http.StatusOK, prefs)
	}
}

func GetCities() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		cities := model.GetCityAll()
		return c.JSON(http.StatusOK, cities)
	}
}

func GetCitiesByPrefCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		pref_code := c.Param("pref_code")
		cities := model.GetCityByPrefCode(pref_code)
		return c.JSON(http.StatusOK, cities)
	}
}

// CityCodeを用いてFrom(現在地)の昨日の天気、今日の予想天気を取得
func GetWeatherFromByCityCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		city_code := c.Param("city_code")
		// cityのLon,Lat等を取得
		city := model.GetLocationByCityCode(city_code)

		//// 並行処理を行う
		var wg sync.WaitGroup
		wg.Add(2)
		var weatherInfoYesterday WeatherInfo
		var weatherInfosToday WeatherInfos
		var err1 error
		var err2 error
		// 昨日の天気を取得
		go func() {
			defer wg.Done()
			weatherInfoYesterday, err1 = GetWeatherYesterday(city.CityLon, city.CityLat)
		}()
		// 1日分(今日)の天気予報を取得
		go func() {
			defer wg.Done()
			weatherInfosToday, err2 = GetWeatherForecast(city.CityLon, city.CityLat, 1)
		}()
		// 並行処理待ち合わせ
		wg.Wait()

		if err1 != nil || err2 != nil {
			return c.JSON(http.StatusServiceUnavailable, "ServiceUnavailable")
		}

		// 昨日と今日をスライスにまとめる
		var weatherInfos WeatherInfos
		weatherInfos = append(weatherInfos, weatherInfoYesterday)
		weatherInfos = append(weatherInfos, weatherInfosToday[0])

		return c.JSON(http.StatusOK, weatherInfos)
	}
}

// CityCodeを用いてTo(目的地)の予想天気を取得
func GetWeatherToByCityCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		city_code := c.Param("city_code")
		// cityのLon,Lat等を取得
		city := model.GetLocationByCityCode(city_code)
		// 8日分(今日+7日間)の天気予報を取得
		weatherInfos, err := GetWeatherForecast(city.CityLon, city.CityLat, 8)
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, "ServiceUnavailable")
		}
		return c.JSON(http.StatusOK, weatherInfos)
	}
}

// 受け取ったJSON(fromCityCode, toCityCode)を用いてfavoritesテーブルに対しINSする(同一レコードが存在する場合は更新日時のみUPD)
func CreateFavoriteFromJson() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		// 受け取ったjsonを構造体にバインド
		favoriteRequired := FavoriteRequired{}
		err := c.Bind(&favoriteRequired)
		if err != nil {
			utils.OutErrorLog("failed to c.Bind", err)
			return c.JSON(http.StatusInternalServerError, "InternalServerError")
		}

		// favoritesテーブルに対しINS(またはUPD)
		resultCode, err := model.CreateFavoriteTransaction(favoriteRequired.Nickname, favoriteRequired.FromCityCode, favoriteRequired.ToCityCode)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "InternalServerError")
		}
		favoriteInsResult := FavoriteInsResult{ResultCode: resultCode}

		return c.JSON(http.StatusCreated, favoriteInsResult)
	}
}

// SelectedFavorites構造体を受け取り、UpdatedAtをtime型からstring型へ変換した構造体を返す関数
func ChangeSelectedFavriteToStrs(selectedFavorites model.SelectedFavorites) SelectedFavoriteStrs {
	// 更新日時を文字列へ変換し、リターン用構造体に代入
	var selectedFavoriteStrs SelectedFavoriteStrs
	for _, v := range selectedFavorites {
		timeDayStr := v.UpdatedAt.Format(config.WEATHER_DATE_FORMAT)
		selectedFavoriteStr := SelectedFavoriteStr{
			Nickname:     v.Nickname,
			FromPrefCode: v.FromPrefCode,
			FromCityCode: v.FromCityCode,
			ToPrefCode:   v.ToPrefCode,
			ToCityCode:   v.ToCityCode,
			FromPrefName: v.FromPrefName,
			FromCityName: v.FromCityName,
			ToPrefName:   v.ToPrefName,
			ToCityName:   v.ToCityName,
			UpdatedAt:    timeDayStr,
		}
		selectedFavoriteStrs = append(selectedFavoriteStrs, selectedFavoriteStr)
	}
	return selectedFavoriteStrs
}

// favoritesテーブルから更新日時降順で全件を取得して返す
func GetFavoriteAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		// DB検索
		selectedFavorites := model.GetFavoriteAll()
		// 更新日時を文字列へ変換
		selectedFavoriteStrs := ChangeSelectedFavriteToStrs(selectedFavorites)

		return c.JSON(http.StatusOK, selectedFavoriteStrs)
	}
}

// favoritesテーブルから更新日時降順でn件までを取得して返す
func GetFavoriteN() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])
		nString := c.Param("n")
		// nがint型へ変換できない場合 or nが0以下の場合、StatusBadRequestでリターン
		n, err := strconv.Atoi(nString)
		if err != nil || n <= 0 {
			return c.JSON(http.StatusBadRequest, "invalid request:"+nString)
		}

		// DB検索(n件まで)
		selectedFavorites := model.GetFavoriteN(n)
		// 更新日時を文字列へ変換
		selectedFavoriteStrs := ChangeSelectedFavriteToStrs(selectedFavorites)

		return c.JSON(http.StatusOK, selectedFavoriteStrs)
	}
}

// favoritesテーブルから更新日時降順でn件までを取得して返す
func GetFavoriteByNickname() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])
		nickname := c.Param("nickname")

		// DB検索
		selectedFavorites := model.GetFavoriteByNickname(nickname)
		// 更新日時を文字列へ変換
		selectedFavoriteStrs := ChangeSelectedFavriteToStrs(selectedFavorites)

		return c.JSON(http.StatusOK, selectedFavoriteStrs)
	}
}

// favoritesテーブルから更新日時降順でn件までを取得して返す
func GetNicknameDistinct() echo.HandlerFunc {
	return func(c echo.Context) error {
		// user access log
		utils.OutInfoLogUserAccess("START", c.RealIP(), c.Request().Header["User-Agent"][0])

		// DB検索
		selectedNicknames := model.GetNicknameDistinct()

		return c.JSON(http.StatusOK, selectedNicknames)
	}
}

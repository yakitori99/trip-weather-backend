package api

import (
	"net/http"
	"sync"
	"time"

	echo "github.com/labstack/echo/v4"
	// 自作パッケージをインポート
	"trip-weather-backend/config"
	"trip-weather-backend/model"
)

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	}
}

func HelloUsername() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		return c.String(http.StatusOK, "Hello World! Your name is "+username)
	}
}

// today-1 ~ +7のdatetimeをstrにしたデータを返す
func GetDatetimes() echo.HandlerFunc {
	return func(c echo.Context) error {
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
		prefs := model.GetPrefAll()
		return c.JSON(http.StatusOK, prefs)
	}
}

func GetCities() echo.HandlerFunc {
	return func(c echo.Context) error {
		cities := model.GetCityAll()
		return c.JSON(http.StatusOK, cities)
	}
}

func GetCitiesByPrefCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		pref_code := c.Param("pref_code")
		cities := model.GetCityByPrefCode(pref_code)
		return c.JSON(http.StatusOK, cities)
	}
}

// CityCodeを用いてFrom(現在地)の昨日の天気、今日の予想天気を取得
func GetWeatherFromByCityCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		city_code := c.Param("city_code")
		// cityのLon,Lat等を取得
		city := model.GetLocationByCityCode(city_code)

		//// 並行処理を行う
		var wg sync.WaitGroup
		wg.Add(2)
		var weatherInfoYesterday WeatherInfo
		var weatherInfosToday WeatherInfos
		// 昨日の天気を取得
		go func() {
			defer wg.Done()
			weatherInfoYesterday = GetWeatherYesterday(city.CityLon, city.CityLat)
		}()
		// 1日分(今日)の天気予報を取得
		go func() {
			defer wg.Done()
			weatherInfosToday = GetWeatherForcast(city.CityLon, city.CityLat, 1)
		}()
		// 並行処理待ち合わせ
		wg.Wait()

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
		city_code := c.Param("city_code")
		// cityのLon,Lat等を取得
		city := model.GetLocationByCityCode(city_code)
		// 8日分(今日+7日間)の天気予報を取得
		weatherInfos := GetWeatherForcast(city.CityLon, city.CityLat, 8)
		return c.JSON(http.StatusOK, weatherInfos)
	}
}

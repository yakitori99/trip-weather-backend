package api_test

import (
	"testing"

	// 自作パッケージをインポート
	"trip-weather-backend/api"
)

func Test_GetDailyWeatherCodeFromHourly(t *testing.T) {
	// 2次元スライスの初期化 //テスト用のweatherCodes
	weatherCodesSlice := [][]int{
		{600, 601, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 201, 211}, // Snow:600
		{501, 501, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 201, 211}, // Thunderstorm:200
		{600, 201, 500, 501, 502, 800, 800, 800, 800, 800, 800, 800, 800, 800, 800}, // Rain:501
		{600, 201, 500, 501, 801, 801, 801, 801, 801, 801, 800, 800, 800, 800, 800}, // Clear:800
		{600, 201, 500, 501, 801, 801, 801, 801, 801, 802, 800, 800, 800, 800, 802}, // few clouds:801
		{500, 500, 801, 801, 801, 801, 800, 800, 800, 800, 802, 802, 802, 802, 802}, // scattered clouds:802
		{500, 500, 801, 801, 801, 801, 800, 800, 800, 800, 802, 802, 802, 802, 803}, // clouds:803
		{701, 701, 701, 701, 701, 701, 701, 701, 701, 701, 701, 701, 701, 701, 701}, // clouds:803
	}
	expectedSlice := []int{
		600,
		200,
		501,
		800,
		801,
		802,
		803,
		803,
	}

	for i := range weatherCodesSlice {
		actual := api.GetDailyWeatherCodeFromHourly(weatherCodesSlice[i])
		expected := expectedSlice[i]
		if actual != expected {
			t.Errorf("got:%v, want:%v", actual, expected)
		} else {
			t.Logf("OK weatherCode:%v", actual)
		}
	}
}

func Test_GetWeatherYesterday(t *testing.T) {
	// 東京の緯度経度
	var lon float64 = 139.691711
	var lat float64 = 35.689499

	weatherInfoYesterday, err := api.GetWeatherYesterday(lon, lat)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("MaxTemp:%v", weatherInfoYesterday.MaxTemp)
		t.Logf("MinTemp:%v", weatherInfoYesterday.MinTemp)
		t.Logf("WeatherCode:%v", weatherInfoYesterday.WeatherCode)
		t.Logf("DateTimeStr:%v", weatherInfoYesterday.DateTimeStr)
	}
}

func Test_GetWeatherForecast(t *testing.T) {
	// 東京の緯度経度
	var lon float64 = 139.691711
	var lat float64 = 35.689499

	getDayNums := []int{1, 5, 7}

	for _, getDayNum := range getDayNums {
		weatherInfos, err := api.GetWeatherForecast(lon, lat, getDayNum)
		if err != nil {
			t.Error(err)
		}
		if len(weatherInfos) != getDayNum {
			t.Errorf("got len:%v, want:%v", len(weatherInfos), getDayNum)
		} else {
			t.Logf("getDayNum:%v", getDayNum)
		}

	}
}

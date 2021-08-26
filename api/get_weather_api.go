package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	// jsonの解析用
	simplejson "github.com/bitly/go-simplejson"

	"trip-weather-backend/config"
)

// 戻り値を整理して保持する構造体
type WeatherInfo struct {
	DateTimeStr string
	MaxTemp     float64
	MinTemp     float64
	WeatherCode int
}
type WeatherInfos []WeatherInfo

// 条件とスライス(どちらもint)から、条件に一致する要素数を返す関数
func CountSpecificNumFromSlise(checkCode int, checkSlice []int) int {
	num := 0
	for _, v := range checkSlice {
		if checkCode == v {
			num += 1
		}
	}
	return num
}

// longitude(経度), latitude(緯度)を受け取り、APIに問い合わせて昨日の天気情報を返す関数
func GetWeatherYesterday(lon float64, lat float64) WeatherInfo {
	fmt.Println("GetWeatherYesterday START", time.Now())
	// Requestインスタンス生成
	request, err := http.NewRequest("GET", config.WEATHER_API_TIMEMACHINE_BASEURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	// グリニッジ標準時における昨日を表すunixTimeを生成
	timeYesterday := time.Now().AddDate(0, 0, -1).Add(time.Hour * time.Duration(9))
	var unixTimeYesterday int64 = timeYesterday.Unix()

	// クエリパラメータ作成
	params := request.URL.Query()
	params.Add("lon", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("lat", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("appid", config.API_KEY)
	params.Add("units", config.WEATHER_API_UNITS)
	params.Add("lang", config.WEATHER_API_LANG)
	params.Add("dt", strconv.FormatInt(unixTimeYesterday, 10))
	request.URL.RawQuery = params.Encode()

	// タイムアウトまでの時間を設定
	timeout := time.Duration(time.Duration(config.WEATHER_API_TIMEOUT_SEC) * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	// for Debug
	fmt.Println(request.URL.String())

	// HTTPリクエスト実行
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	// 関数を抜ける際に必ずresponseをcloseする
	defer response.Body.Close()

	// レスポンスを取得
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// レスポンスのjsonを解析 go-simplejsonを利用
	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatal(err)
	}

	//// hourlyのデータからdailyの情報を判定
	// 日時(時間は使わないので、最初の要素だけ取得)
	dateTime := time.Unix(js.Get("hourly").GetIndex(0).Get("dt").MustInt64(), 0)
	dateTimeStr := dateTime.Format(config.WEATHER_DATE_FORMAT)

	var tempMax float64 = -100.0
	var tempMin float64 = 100.0
	var weatherCodes []int
	for i, _ := range js.Get("hourly").MustArray() {
		// 気温
		temp := js.Get("hourly").GetIndex(i).Get("temp").MustFloat64()
		if temp > tempMax {
			tempMax = temp
		}
		if temp < tempMin {
			tempMin = temp
		}
		// 天気コードを天気判定用スライスへ追加
		if len(weatherCodes) < config.WEATHER_CHECK_LENGTH {
			_weatherCode := js.Get("hourly").GetIndex(i).Get("weather").GetIndex(0).Get("id").MustInt()
			weatherCodes = append(weatherCodes, _weatherCode)
		}
	}
	//// dailyで見た天気を判定
	// 1桁判定用スライス
	var weatherCodesOne []int
	for _, _weatherCode := range weatherCodes {
		weatherCodesOne = append(weatherCodesOne, _weatherCode/100)
	}
	// 判定
	weatherCode := -1
	for _, row := range config.WEATHER_CODE_THRESHS {
		checkCode := row[0]
		outCode := row[1]
		thresh := row[2]
		var num int
		if checkCode <= 9 { //1桁判定
			num = CountSpecificNumFromSlise(checkCode, weatherCodesOne)
		} else { //3桁判定
			num = CountSpecificNumFromSlise(checkCode, weatherCodes)
		}

		if num >= thresh {
			weatherCode = outCode
			break
		}
	}
	// ここまでで天気が決まらない場合、曇りとみなす
	if weatherCode == -1 {
		weatherCode = config.WEATHER_CODE_CLOUDS
	}

	//// 戻り値構造体に代入
	weatherInfo := WeatherInfo{dateTimeStr, tempMax, tempMin, weatherCode}
	fmt.Println("GetWeatherYesterday END", time.Now())
	return weatherInfo
}

// longitude(経度), latitude(緯度), 天気予報取得日数を受け取り、APIに問い合わせて天気予報を返す関数
func GetWeatherForcast(lon float64, lat float64, getDayNum int) WeatherInfos {
	fmt.Println("GetWeatherForcast START", time.Now())
	// Requestインスタンス生成
	request, err := http.NewRequest("GET", config.WEATHER_API_FORECAST_BASEURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	// クエリパラメータ作成
	params := request.URL.Query()
	params.Add("lon", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("lat", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("appid", config.API_KEY)
	params.Add("units", config.WEATHER_API_UNITS)
	params.Add("lang", config.WEATHER_API_LANG)
	request.URL.RawQuery = params.Encode()

	// タイムアウトまでの時間を設定
	timeout := time.Duration(time.Duration(config.WEATHER_API_TIMEOUT_SEC) * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	// for Debug
	fmt.Println(request.URL.String())

	// HTTPリクエスト実行
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	// 関数を抜ける際に必ずresponseをcloseする
	defer response.Body.Close()

	// レスポンスを取得
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// レスポンスのjsonを解析 go-simplejsonを利用
	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatal(err)
	}
	// レスポンスを戻り値構造体に代入
	var weatherInfos WeatherInfos
	for i := 0; i < getDayNum; i++ {
		// 要素を取得し型変換
		dateTime := time.Unix(js.Get("daily").GetIndex(i).Get("dt").MustInt64(), 0)
		dateTimeStr := dateTime.Format(config.WEATHER_DATE_FORMAT)
		tempMax := js.Get("daily").GetIndex(i).Get("temp").Get("max").MustFloat64()
		tempMin := js.Get("daily").GetIndex(i).Get("temp").Get("min").MustFloat64()
		weatherCode := js.Get("daily").GetIndex(i).Get("weather").GetIndex(0).Get("id").MustInt()

		// 戻り値用sliceに追加
		weatherInfos = append(weatherInfos, WeatherInfo{dateTimeStr, tempMax, tempMin, weatherCode})
	}
	fmt.Println("GetWeatherForcast END", time.Now())
	return weatherInfos
}

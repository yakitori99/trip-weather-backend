package config

// 定数
const (
	// API Server
	// PORT string = "" // 使いたいポート番号を指定

	// Weather API
	// API_KEY                         string = "" // 有効なAPI_KEYを指定
	WEATHER_API_TIMEMACHINE_BASEURL string = "https://api.openweathermap.org/data/2.5/onecall/timemachine"
	WEATHER_API_FORECAST_BASEURL    string = "https://api.openweathermap.org/data/2.5/onecall"
	WEATHER_API_UNITS               string = "metric"
	WEATHER_API_LANG                string = "ja"
	WEATHER_API_TIMEOUT_SEC         int    = 5
	// Weather関連
	// 0埋めあり年月日, 24h表現 -- 具体的な数字だが、yyyy/MM/DDのようなお決まりの表現。
	WEATHER_DATE_FORMAT string = "2006/01/02 15:04:05"
	// 9-23時
	WEATHER_CHECK_LENGTH int = 15
	// 普通の曇り
	WEATHER_CODE_CLOUDS int = 803

	//// DB
	DB_PATH string = "db/trip_weather.db"
	// execute code
	DONE_ERR int = -1
	DONE_INS int = 1
	DONE_UPD int = 2
)

//// Goではarray,slice,mapは定数として扱えないため別途宣言
// hourlyからdailyのweatherCodeを判定するための配列
// weatherCode と 判定するしきい値のセット。hourlyのweatherCodeを上から順に評価し、条件を満たしたらそこで終了
// 以下のチェック用コードが1桁の場合、例えば条件が5なら、Codeが5XXのいずれでも条件を満たすとみなす
var WEATHER_CODE_THRESHS [6][3]int = [6][3]int{
	// チェック用、output用Code, 個数のしきい値
	{6, 600, 2},   // Snow
	{2, 200, 2},   // Thunderstorm
	{5, 501, 3},   // Rain
	{800, 800, 5}, // Clear //晴れ
	{801, 801, 5}, // few clouds // 晴れ時々曇り
	{802, 802, 5}, // scattered clouds // 晴れ時々曇り
}

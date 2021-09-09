package route

import (
	// 自作パッケージをインポート
	"trip-weather-backend/api"

	// インポートし、echoという別名で利用する
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute() *echo.Echo {
	// echoのインスタンス作成
	e := echo.New()

	// CORS(Cross Origin Resource Sharing)対応 -- vue.jsからの呼び出しを可能とする
	e.Use(middleware.CORS())

	// API起動確認用
	e.GET("/hello", api.Hello())
	e.GET("/hello/:username", api.HelloUsername())

	// today-1 ~ +7のdatetimeをstrにしたデータを返す
	// クライアントサイドの時間は信用できない可能性があるため、サーバサイドで時間を取得する
	e.GET("/get_datetimes", api.GetDatetimes())

	// pref一覧を返す
	e.GET("/get_prefs", api.GetPrefs())

	// city一覧を返す
	e.GET("/get_cities", api.GetCities())
	// pref_codeで検索しcity一覧を返す
	e.GET("/get_cities/:pref_code", api.GetCitiesByPrefCode())

	// city_codeを用いて、現在地のWeatherHistory, WeatherForcast情報を返す
	e.GET("/get_weather_from/:city_code", api.GetWeatherFromByCityCode())
	// city_codeを用いて、目的地のWeatherForcast情報を返す
	e.GET("/get_weather_to/:city_code", api.GetWeatherToByCityCode())

	// jsonを受け取り、favoritesテーブルに対しINSまたはUPDする
	e.POST("/favorites", api.CreateFavoriteFromJson())

	return e
}

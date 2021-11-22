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
	e.GET("/datetimes", api.GetDatetimes())

	// pref一覧を返す
	e.GET("/prefs", api.GetPrefs())

	// city一覧を返す
	e.GET("/cities", api.GetCities())
	// pref_codeで検索しcity一覧を返す
	e.GET("/cities/by/:pref_code", api.GetCitiesByPrefCode())

	// city_codeを用いて、現在地のWeatherHistory, WeatherForcast情報を返す
	e.GET("/weather_from/:city_code", api.GetWeatherFromByCityCode())
	// city_codeを用いて、目的地のWeatherForcast情報を返す
	e.GET("/weather_to/:city_code", api.GetWeatherToByCityCode())

	// jsonを受け取り、favoritesテーブルに対しINSまたはUPDする
	e.POST("/favorites", api.CreateFavoriteFromJson())

	// favoritesテーブルから更新日時降順で全件を取得して返す
	e.GET("/favorites", api.GetFavoriteAll())
	// favoritesテーブルから更新日時降順でn件までを取得して返す
	e.GET("/favorites/to/:n", api.GetFavoriteN())
	// favoritesテーブルからニックネームをキーに検索し、更新日時降順で全件を取得して返す
	e.GET("/favorites/by/:nickname", api.GetFavoriteByNickname())

	// favoritesテーブルから重複しないニックネーム一覧を取得し、ニックネームの昇順で返す
	e.GET("/nicknames", api.GetNicknameDistinct())

	return e
}

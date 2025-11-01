package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherReading struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int64 `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

var exampleWeather = WeatherReading{
	Coord: struct {
		Lon float64 "json:\"lon\""
		Lat float64 "json:\"lat\""
	}{
		Lon: 7.367,
		Lat: 45.133,
	},
	Weather: []struct {
		ID          int    "json:\"id\""
		Main        string "json:\"main\""
		Description string "json:\"description\""
		Icon        string "json:\"icon\""
	}{
		{ID: 501, Main: "Rain", Description: "moderate rain", Icon: "10d"},
	},
	Base: "stations",
	Main: struct {
		Temp      float64 "json:\"temp\""
		FeelsLike float64 "json:\"feels_like\""
		TempMin   float64 "json:\"temp_min\""
		TempMax   float64 "json:\"temp_max\""
		Pressure  int     "json:\"pressure\""
		Humidity  int     "json:\"humidity\""
		SeaLevel  int     "json:\"sea_level\""
		GrndLevel int     "json:\"grnd_level\""
	}{
		Temp: 284.2, FeelsLike: 282.93, TempMin: 283.06, TempMax: 286.82,
		Pressure: 1021, Humidity: 60, SeaLevel: 1021, GrndLevel: 910,
	},
	Visibility: 10000,
	Wind: struct {
		Speed float64 "json:\"speed\""
		Deg   int     "json:\"deg\""
		Gust  float64 "json:\"gust\""
	}{
		Speed: 4.09, Deg: 121, Gust: 3.47,
	},
	Rain: struct {
		OneH float64 "json:\"1h\""
	}{
		OneH: 2.73,
	},
	Clouds: struct {
		All int "json:\"all\""
	}{
		All: 83,
	},
	Dt: 1726660758,
	Sys: struct {
		Type    int    "json:\"type\""
		ID      int    "json:\"id\""
		Country string "json:\"country\""
		Sunrise int64  "json:\"sunrise\""
		Sunset  int64  "json:\"sunset\""
	}{
		Type: 1, ID: 6736, Country: "IT", Sunrise: 1726636384, Sunset: 1726680975,
	},
	Timezone: 7200,
	ID:       3165523,
	Name:     "Province of Turin",
	Cod:      200,
}

// GET endpoint
func getWeather(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, exampleWeather)
}

func main() {
	router := gin.Default()
	router.GET("/weather", getWeather)
	router.Run(":8080")
}

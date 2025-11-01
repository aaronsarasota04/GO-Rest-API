package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

var exampleWeather []WeatherReading

func init() {
	jsonData := `
	[
	  {
	    "coord": {"lon": 7.367, "lat": 45.133},
	    "weather": [{"id": 501, "main": "Rain", "description": "moderate rain", "icon": "10d"}],
	    "base": "stations",
	    "main": {"temp": 284.2, "feels_like": 282.93, "temp_min": 283.06, "temp_max": 286.82, "pressure": 1021, "humidity": 60, "sea_level": 1021, "grnd_level": 910},
	    "visibility": 10000,
	    "wind": {"speed": 4.09, "deg": 121, "gust": 3.47},
	    "rain": {"1h": 2.73},
	    "clouds": {"all": 83},
	    "dt": 1726660758,
	    "sys": {"type": 1, "id": 6736, "country": "IT", "sunrise": 1726636384, "sunset": 1726680975},
	    "timezone": 7200,
	    "id": 3165523,
	    "name": "Province of Turin",
	    "cod": 200
	  }
	]`

	json.Unmarshal([]byte(jsonData), &exampleWeather)
}

func addWeather(c *gin.Context) {
	var newdata WeatherReading

	if err := c.BindJSON(&newdata); err != nil {
		return
	}
	exampleWeather = append(exampleWeather, newdata)
	c.IndentedJSON(http.StatusCreated, newdata)
}

// PUT /weather/:id - Update an existing weather record
func updateWeather(c *gin.Context) {
	id := c.Param("id") // get the ID from URL
	var updatedData WeatherReading

	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	for i, w := range exampleWeather {
		if fmt.Sprintf("%d", w.ID) == id {
			exampleWeather[i] = updatedData
			c.IndentedJSON(http.StatusOK, updatedData)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
}

// DELETE /weather/:id - Remove a weather record
func deleteWeather(c *gin.Context) {
	id := c.Param("id") // get the ID from URL

	for i, w := range exampleWeather {
		if fmt.Sprintf("%d", w.ID) == id {
			// remove the element from slice
			exampleWeather = append(exampleWeather[:i], exampleWeather[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
}

// get an existing weather record
func getWeatherByCoordinate(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	if latStr == "" || lonStr == "" {
		c.IndentedJSON(http.StatusOK, exampleWeather)
		return
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude or longitude"})
		return
	}

	for _, w := range exampleWeather {
		if w.Coord.Lat == lat && w.Coord.Lon == lon {
			c.IndentedJSON(http.StatusOK, w)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
}

func main() {
	router := gin.Default()
	router.GET("/weather", getWeatherByCoordinate)
	router.POST("/weather", addWeather)
	router.PUT("/weather/:id", updateWeather)
	router.DELETE("/weather/:id", deleteWeather)
	router.Run(":8080")
}

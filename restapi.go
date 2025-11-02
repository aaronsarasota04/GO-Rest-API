package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WeatherReading struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Elevation            float64 `json:"elevation"`
	GenerationTimeMS     float64 `json:"generationtime_ms"`
	UTCOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`

	Hourly struct {
		Time          []string  `json:"time"`                 // ISO8601 timestamps
		Temperature2M []float64 `json:"temperature_2m"`       // °C
		WindSpeed10M  []float64 `json:"wind_speed_10m"`       // km/h
		Precipitation []float64 `json:"precipitation"`        // mm
		PrecipType    []int     `json:"weather_code"`         // 0=clear, 1-3=cloudy, 51-67=rain, 71-77=snow, etc.
		Humidity2M    []float64 `json:"relative_humidity_2m"` // %
	} `json:"hourly"`

	HourlyUnits struct {
		Temperature2M string `json:"temperature_2m"`       // °C
		WindSpeed10M  string `json:"wind_speed_10m"`       // km/h
		Precipitation string `json:"precipitation"`        // mm
		PrecipType    string `json:"weather_code"`         // code
		Humidity2M    string `json:"relative_humidity_2m"` // %
	} `json:"hourly_units"`
}

var exampleWeather []WeatherReading

func init() {
	jsonData := `
	[
	  {
  "latitude": 52.52,
  "longitude": 13.41,
  "elevation": 44.8,
  "generationtime_ms": 2.5,
  "utc_offset_seconds": 0,
  "timezone": "Europe/Berlin",
  "timezone_abbreviation": "CET",
  "hourly": {
    "time": ["2025-11-01T00:00", "2025-11-01T01:00", "2025-11-01T02:00"],
    "temperature_2m": [8.1, 7.8, 7.5],
    "wind_speed_10m": [12.3, 11.8, 11.5],
    "relative_humidity_2m": [80, 82, 85],
    "precipitation": [0.0, 0.2, 0.0],
    "precipitation_type": ["none", "rain", "none"]
  },
  "hourly_units": {
    "temperature_2m": "°C",
    "wind_speed_10m": "km/h",
    "relative_humidity_2m": "%",
    "precipitation": "mm",
    "precipitation_type": "string"
  }
}
	]`

	json.Unmarshal([]byte(jsonData), &exampleWeather)
}

// DELETE /weather/:id - Remove a weather record
func deleteWeather(c *gin.Context) {
	latStr := c.Query("Latitude")
	lonStr := c.Query("Longitude")

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude or longitude"})
		return
	}

	for i, w := range exampleWeather {
		if w.Latitude == lat && w.Longitude == lon {
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
	latStr := c.Query("Latitude")
	lonStr := c.Query("Longitude")

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
		if w.Latitude == lat && w.Longitude == lon {
			c.IndentedJSON(http.StatusOK, w)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
}

func fetchWeather(c *gin.Context) {
	latStr := c.Query("Latitude")
	lonStr := c.Query("Longitude")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and Longitude are required"})
		return
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude or longitude"})
		return
	}

	// Construct Open-Meteo URL with desired hourly variables
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m,precipitation,weather_code&timezone=auto",
		lat, lon,
	)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather"})
		return
	}
	defer resp.Body.Close()

	var data WeatherReading
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse weather data"})
		return
	}

	// Optionally, append or replace your local slice
	exampleWeather = append(exampleWeather, data)

	c.IndentedJSON(http.StatusOK, data)
}

func main() {
	router := gin.Default()
	router.GET("/weather", getWeatherByCoordinate)
	router.POST("/weather/fetch", fetchWeather)
	router.DELETE("/weather", deleteWeather)
	router.Run(":8080")
}

package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	WEATHER_BASE_URL = "http://api.openweathermap.org/data/2.5/weather"
	KELVIN_CONSTANT  = 273.15
)

type WeatherResponse struct {
	Main struct {
		FeelsLike float64 `json:"feels_like"`
		Temp      float64 `json:"temp"`
	} `json:"main"`
	Name    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

type WeatherService struct {
	apiKey string
}

// NewWeatherService initialises a WeatherService that is ready to use.
func NewWeatherService() *WeatherService {
	key, ok := os.LookupEnv("WEATHER_API_KEY")
	if !ok {
		log.Fatal("WEATHER_API_KEY must be set")
	}

	return &WeatherService{
		apiKey: key,
	}
}

// GetData gets the weather data for the given city
func (w *WeatherService) GetData(city string) (*WeatherResponse, error) {
	weatherURL := w.getURL(city)
	response, err := http.Get(weatherURL)
	if err != nil {
		return nil, err
	}
	jsonResponse := new(WeatherResponse)
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		return nil, err
	}

	return jsonResponse, nil
}

// GetWeatherURL constructs the weather url based on city and API key
func (w *WeatherService) getURL(city string) string {
	return fmt.Sprintf("%s?q=%s&appid=%s", WEATHER_BASE_URL, city, w.apiKey)
}

// ConvertCelsius takes a Kelvin temp and returns the Celsius number
func ConvertCelsius(kTemp float64) string {
	return fmt.Sprintf("%.2f", kTemp-KELVIN_CONSTANT)
}

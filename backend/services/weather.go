package services

import (
	"encoding/json"
	"fmt"
	"forecaster-go-react/models"
	"net/http"
)

const weatherAPIBaseURL = "http://api.openweathermap.org/data/2.5"
const apiKey = "your_api_key_here"

func FetchCurrentWeather(city string) (*models.Weather, error) {
	url := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=metric", weatherAPIBaseURL, city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather data: %s", resp.Status)
	}

	var weatherData models.Weather
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}

package services

import (
	"encoding/json"
	"fmt"
	"forecaster-go-react/models"
	"net/http"
)

func FetchWeatherForecast(city string) (*models.Forecast, error) {
	url := fmt.Sprintf("%s/forecast/daily?q=%s&appid=%s&units=metric&cnt=7", weatherAPIBaseURL, city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather forecast: %s", resp.Status)
	}

	var forecastData models.Forecast
	if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
		return nil, err
	}

	return &forecastData, nil
}

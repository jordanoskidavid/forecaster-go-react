package handlers

import (
	"encoding/json"
	"fmt"
	"forecaster-go-react/models"
	"io"
	"log"
	"net/http"
)

func GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	country := r.URL.Query().Get("country")

	if city == "" || country == "" {
		http.Error(w, "City and Country parameters are required", http.StatusBadRequest)
		return
	}

	geoAPI := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s", city)

	resp, err := http.Get(geoAPI)
	if err != nil {
		http.Error(w, "Failed to fetch geocode data", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	var geoResponse models.GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResponse); err != nil {
		http.Error(w, "Failed to parse geocode response", http.StatusInternalServerError)
		return
	}
	var selectedResult *models.GeocodeResult
	for _, result := range geoResponse.Results {
		if result.Country == country {
			selectedResult = &result
			break
		}
	}
	if selectedResult == nil {
		http.Error(w, "No matching city found in the specified country", http.StatusNotFound)
		return
	}
	weatherAPI := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true",
		selectedResult.Latitude, selectedResult.Longitude,
	)
	weatherResp, err := http.Get(weatherAPI)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(weatherResp.Body)

	body, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		http.Error(w, "Failed to read weather response body", http.StatusInternalServerError)
		return
	}

	var weather models.WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		http.Error(w, "Failed to parse weather response", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"city":    selectedResult.Name,
		"country": selectedResult.Country,
		"weather": weather.CurrentWeather,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"forecaster-go-react/models"
	"io"
	"log"
	"net/http"
	"time"
)

func GetHourlyForecast(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	country := r.URL.Query().Get("country")
	if city == "" || country == "" {
		http.Error(w, "City and Country parameters are required", http.StatusBadRequest)
		return
	}

	client := &http.Client{
		Timeout: 120 * time.Second, // If sometimes makes problem, here you can increase the timeout
	}

	geoAPI := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s", city)
	resp, err := client.Get(geoAPI)
	if err != nil {
		log.Println("Error fetching geocode data:", err)
		http.Error(w, "Failed to fetch geocode data", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Geocoding API returned status: %d", resp.StatusCode)
		http.Error(w, "Failed to fetch geocode data", http.StatusInternalServerError)
		return
	}

	var geoResponse models.GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResponse); err != nil {
		log.Println("Error decoding geocode response:", err)
		http.Error(w, "Failed to parse geocode response", http.StatusInternalServerError)
		return
	}

	if len(geoResponse.Results) == 0 {
		log.Println("No geocode results found")
		http.Error(w, "No results found for the specified city", http.StatusNotFound)
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
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,windspeed_10m&timezone=GMT",
		selectedResult.Latitude, selectedResult.Longitude,
	)

	weatherResp, err := client.Get(weatherAPI)
	if err != nil {
		log.Println("Error fetching weather data:", err)
		http.Error(w, "Failed to fetch hourly forecast data", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body:", err)
		}
	}(weatherResp.Body)

	if weatherResp.StatusCode != http.StatusOK {
		log.Printf("Weather API returned status: %d", weatherResp.StatusCode)
		http.Error(w, "Failed to fetch hourly forecast data", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		log.Println("Error reading weather response body:", err)
		http.Error(w, "Failed to read hourly forecast response", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body:", err)
		}
	}(weatherResp.Body)

	var weather struct {
		Hourly struct {
			Time        []string  `json:"time"`
			Temperature []float64 `json:"temperature_2m"`
			Windspeed   []float64 `json:"windspeed_10m"`
		} `json:"hourly"`
	}
	if err := json.Unmarshal(body, &weather); err != nil {
		log.Println("Error decoding weather response:", err)
		http.Error(w, "Failed to parse hourly forecast response", http.StatusInternalServerError)
		return
	}

	if len(weather.Hourly.Time) == 0 {
		log.Println("Weather API returned incomplete data")
		http.Error(w, "Incomplete hourly forecast data", http.StatusInternalServerError)
		return
	}

	var hourlyForecast []models.HourlyForecast
	for i := range weather.Hourly.Time {
		hourlyForecast = append(hourlyForecast, models.HourlyForecast{
			Time:        weather.Hourly.Time[i],
			Temperature: weather.Hourly.Temperature[i],
			Windspeed:   weather.Hourly.Windspeed[i],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"city":    selectedResult.Name,
		"country": selectedResult.Country,
		"hourly":  hourlyForecast,
	})
	if err != nil {
		return
	}
}

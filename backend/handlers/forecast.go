package handlers

import (
	"encoding/json"
	"forecaster-go-react/services"
	"net/http"
)

func GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	forecastData, err := services.FetchWeatherForecast(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(forecastData)
}

package main

import (
	"forecaster-go-react/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/get-weather", handlers.GetCurrentWeather)
	http.HandleFunc("/api/get-hourly-forecast", handlers.GetHourlyForecast)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

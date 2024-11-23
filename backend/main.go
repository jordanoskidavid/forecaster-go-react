package main

import (
	"forecaster-go-react/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/weather", handlers.GetCurrentWeather)
	http.HandleFunc("/api/forecast", handlers.GetWeatherForecast)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

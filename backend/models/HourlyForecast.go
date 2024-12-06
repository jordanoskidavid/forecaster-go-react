package models

type HourlyForecast struct {
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
	Windspeed   float64 `json:"windspeed"`
}

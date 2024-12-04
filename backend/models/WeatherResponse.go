package models

type CurrentWeather struct {
	Time          string  `json:"time"`
	Temperature   float64 `json:"temperature"`
	Windspeed     float64 `json:"windspeed"`
	WindDirection float64 `json:"winddirection"`
	IsDay         int     `json:"is_day"`
	WeatherCode   int     `json:"weathercode"`
}

type WeatherResponse struct {
	CurrentWeather CurrentWeather `json:"current_weather"`
}

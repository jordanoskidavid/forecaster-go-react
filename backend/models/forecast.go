package models

type Forecast struct {
	City     string `json:"city"`
	Forecast []struct {
		Day         string  `json:"day"`
		Temperature float64 `json:"temperature"`
		Description string  `json:"description"`
	} `json:"forecast"`
}

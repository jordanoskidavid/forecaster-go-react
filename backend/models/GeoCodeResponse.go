package models

type GeocodeResponse struct {
	Results []GeocodeResult `json:"results"`
}

type GeocodeResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
}

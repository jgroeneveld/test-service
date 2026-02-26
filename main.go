package main

import (
	"encoding/json"
	"net/http"
)

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature_c"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity_pct"`
	WindSpeed   float64 `json:"wind_speed_kmh"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Weather{
		City:        "Amsterdam",
		Temperature: 7.4,
		Condition:   "Partly cloudy",
		Humidity:    78,
		WindSpeed:   19.2,
	})
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

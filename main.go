package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature_c"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity_pct"`
	WindSpeed   float64 `json:"wind_speed_kmh"`
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
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
	http.HandleFunc("/", loggingMiddleware(indexHandler))
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Weather struct {
	ID          int     `json:"id"`
	City        string  `json:"city"`
	Temperature float64 `json:"temperature_c"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity_pct"`
	WindSpeed   float64 `json:"wind_speed_kmh"`
}

var db *sql.DB

func migrate() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS weather (
		id          SERIAL PRIMARY KEY,
		city        TEXT NOT NULL,
		temperature NUMERIC(5,2) NOT NULL,
		condition   TEXT NOT NULL,
		humidity    INT NOT NULL,
		wind_speed  NUMERIC(5,2) NOT NULL
	)`)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}
}

func seed() {
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM weather`).Scan(&count)
	if count > 0 {
		return
	}

	rows := []Weather{
		{City: "Amsterdam", Temperature: 7.4, Condition: "Partly cloudy", Humidity: 78, WindSpeed: 19.2},
		{City: "London", Temperature: 9.1, Condition: "Overcast", Humidity: 82, WindSpeed: 24.5},
		{City: "Berlin", Temperature: 3.2, Condition: "Snowing", Humidity: 90, WindSpeed: 11.0},
		{City: "Paris", Temperature: 11.0, Condition: "Sunny", Humidity: 65, WindSpeed: 15.3},
		{City: "Madrid", Temperature: 16.5, Condition: "Clear", Humidity: 45, WindSpeed: 8.7},
		{City: "Rome", Temperature: 14.2, Condition: "Sunny", Humidity: 55, WindSpeed: 10.1},
		{City: "Vienna", Temperature: 4.8, Condition: "Foggy", Humidity: 88, WindSpeed: 6.4},
		{City: "Zurich", Temperature: 2.1, Condition: "Snowing", Humidity: 92, WindSpeed: 9.3},
		{City: "Brussels", Temperature: 8.3, Condition: "Rainy", Humidity: 85, WindSpeed: 22.0},
		{City: "Stockholm", Temperature: -1.4, Condition: "Clear", Humidity: 70, WindSpeed: 13.6},
		{City: "Oslo", Temperature: -3.2, Condition: "Snowing", Humidity: 75, WindSpeed: 7.8},
		{City: "Copenhagen", Temperature: 5.5, Condition: "Cloudy", Humidity: 80, WindSpeed: 18.4},
	}

	for _, r := range rows {
		_, err := db.Exec(
			`INSERT INTO weather (city, temperature, condition, humidity, wind_speed) VALUES ($1,$2,$3,$4,$5)`,
			r.City, r.Temperature, r.Condition, r.Humidity, r.WindSpeed,
		)
		if err != nil {
			log.Fatalf("seed: %v", err)
		}
	}
	log.Printf("seeded %d weather rows", len(rows))
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, city, temperature, condition, humidity, wind_speed FROM weather ORDER BY id LIMIT 10`)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("query: %v", err)
		return
	}
	defer rows.Close()

	results := []Weather{}
	for rows.Next() {
		var wr Weather
		if err := rows.Scan(&wr.ID, &wr.City, &wr.Temperature, &wr.Condition, &wr.Humidity, &wr.WindSpeed); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			log.Printf("scan: %v", err)
			return
		}
		results = append(results, wr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}
	log.Println("connected to database")

	migrate()
	seed()

	http.HandleFunc("/", loggingMiddleware(indexHandler))
	http.HandleFunc("/health", loggingMiddleware(healthHandler))
	log.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}

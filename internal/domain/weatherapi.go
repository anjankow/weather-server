package domain

import (
	"encoding/json"
	"time"
)

type WeatherAPIDayForecastQuery struct {
	Longitude float64
	Latitude float64
	Date time.Time
}

type WeatherAPIDayForecast struct {
	// Data in the raw form  (no manipulation requested)
    Data json.RawMessage
}
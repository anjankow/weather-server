package domain

import (
	"encoding/json"
	"time"
)

type Location struct {
	Longitude float64
	Latitude  float64
}

type ForecastQuery struct {
	Location
	FromDay   time.Time
	NumOfDays int
}

type DayForecastQuery struct {
	Location
	Day time.Time
}

type DayForecastRaw json.RawMessage

// ForecastResponse is an array of weather forecasts for n consequitive days.
// Index 0 in the array maps to Day1.
type ForecastResponse []DayForecastRaw

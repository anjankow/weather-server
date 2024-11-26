package domain

import (
	"time"
)

const LongitudeMax float64 = 180
const LongitudeMin float64 = -180
const LatitudeMax float64 = 90
const LatitudeMin float64 = -90

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

type DayForecastRaw map[string]interface{}
type DayForecastSlice []DayForecastRaw

// ForecastResponse is an array of weather forecasts for n consequitive days.
// Index 0 in the array maps to Day1.
type ForecastResponse struct {
	DayForecasts DayForecastSlice
	APIName      string
}

type ForecastAggregate []ForecastResponse

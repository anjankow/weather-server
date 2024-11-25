package app

import (
	openmeteo "weather-server/internal/app/open_meteo"
	weatherapi "weather-server/internal/app/weather_api"
)

type Aggregator struct {
	openMeteo openmeteo.Service
	weatherAPI weatherapi.Service
}

func NewAggregator(openMeteo openmeteo.Service, weatherAPI weatherapi.Service) Aggregator {
	return Aggregator{
		openMeteo: openMeteo,
		weatherAPI: weatherAPI,
	}
}
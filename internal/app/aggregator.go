package app

import (
	openmeteo "weather-server/internal/app/open_meteo"
	"weather-server/internal/app/weatherapi"
)

type Aggregator struct {
	openMeteo  openmeteo.Service
	weatherAPI weatherapi.Service
}

func NewAggregator(openMeteo openmeteo.Service, weatherAPI weatherapi.Service) Aggregator {
	return Aggregator{
		openMeteo:  openMeteo,
		weatherAPI: weatherAPI,
	}
}

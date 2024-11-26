package app

import (
	"fmt"
	"weather-server/internal/domain"
	"weather-server/internal/services"
)

type App struct {
	Aggregator services.Aggregator
}

type Dependencies struct {
	ForecastProviders map[string]domain.Client
}

func New(deps Dependencies) (App,error) {
	aggr,err := services.NewAggregator(deps.ForecastProviders)
	if err != nil {
		return App{}, fmt.Errorf("failed to initialize Aggregator: %w", err)
	}

	return App{
		Aggregator: aggr,
	}, nil
}
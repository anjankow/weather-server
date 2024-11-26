package app

import (
	"errors"
	"fmt"
	"weather-server/internal/domain"
	"weather-server/internal/services"
)

type App struct {
	Aggregator services.Aggregator
}

type Config struct {
	NumOfForecastDays int
}

type Dependencies struct {
	ForecastProviders map[string]domain.Client
}

func New(cfg Config, deps Dependencies) (App, error) {
	if cfg.NumOfForecastDays <= 0 {
		return App{}, errors.New("number of days must be > 0")
	}

	aggr, err := services.NewAggregator(cfg.NumOfForecastDays, deps.ForecastProviders)
	if err != nil {
		return App{}, fmt.Errorf("failed to initialize Aggregator: %w", err)
	}

	return App{
		Aggregator: aggr,
	}, nil
}

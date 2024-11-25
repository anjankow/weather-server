package openmeteo

import (
	"context"
	"weather-server/internal/domain"
)

const ProviderName = "OpenMeteo"

type Service struct {
}

func NewService() Service {
	return Service{}
}

// todo
func (s Service) GetForecast(ctx context.Context, query domain.ForecastQuery) (domain.DayForecastSlice, error) {
	return nil, nil
}

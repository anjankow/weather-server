package weatherapi

import (
	"context"
	"weather-server/internal/domain"
)





type Client interface {
	GetForecast(ctx context.Context, query domain.WeatherAPIDayForecastQuery) (domain.WeatherAPIDayForecast, error)
}


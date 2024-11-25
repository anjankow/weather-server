package weatherapiclient

import (
	"context"
	"weather-server/internal/domain"
)

type Mock struct {

}

func (m *Mock) GetForecast(ctx context.Context, query domain.WeatherAPIDayForecastQuery) (domain.WeatherAPIDayForecast, error) {
	return domain.WeatherAPIDayForecast{}, nil
}
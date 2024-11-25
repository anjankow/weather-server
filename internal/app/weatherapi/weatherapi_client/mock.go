package weatherapiclient

import (
	"context"
	"weather-server/internal/domain"
)

type Mock struct {

}

func (m *Mock) GetForecast(ctx context.Context, query domain.WeatherAPIForecastQuery) (domain.WeatherAPIForecast, error) {
	return domain.WeatherAPIForecast{}, nil
}
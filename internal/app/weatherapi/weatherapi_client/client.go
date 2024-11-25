package weatherapiclient

import (
	"context"
	"weather-server/internal/domain"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) Client{
	return Client{apiKey}
}

func (c Client) GetForecast(ctx context.Context, query domain.WeatherAPIForecastQuery) (domain.WeatherAPIForecast, error) {
	
	return domain.WeatherAPIForecast{}, nil
}
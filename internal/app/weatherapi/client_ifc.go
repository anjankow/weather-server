package weatherapi

import (
	"context"
	"encoding/json"
	"weather-server/internal/domain"
)

type Client interface {
	GetDayForecast(ctx context.Context, query domain.DayForecastQuery) (json.RawMessage, error)
}

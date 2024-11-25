package weatherapiclient

import (
	"context"
	"encoding/json"
	"weather-server/internal/domain"
)

type Mock struct {
}

func (m *Mock) GetDayForecast(ctx context.Context, query domain.DayForecastQuery) (json.RawMessage, error) {
	return json.RawMessage{}, nil
}

package app

import (
	"context"
	"encoding/json"
	"weather-server/internal/domain"
)

//go:generate go run go.uber.org/mock/mockgen -source client_ifc.go -package mock -destination=./mock/client.go
type Client interface {
	GetDayForecast(ctx context.Context, query domain.DayForecastQuery) (json.RawMessage, error)
}

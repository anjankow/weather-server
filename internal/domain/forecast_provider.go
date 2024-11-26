package domain

import (
	"context"
)

//go:generate go run go.uber.org/mock/mockgen -source forecast_provider.go -package mock -destination=./../forecast_providers/mock/client.go
type Client interface {
	GetDayForecast(ctx context.Context, query DayForecastQuery) (DayForecastRaw, error)
}

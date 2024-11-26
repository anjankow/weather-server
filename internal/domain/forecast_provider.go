package domain

import (
	"context"
	"encoding/json"
)

//go:generate go run go.uber.org/mock/mockgen -source client_ifc.go -package mock -destination=./../providers/mock/client.go
type Client interface {
	GetDayForecast(ctx context.Context, query DayForecastQuery) (json.RawMessage, error)
}

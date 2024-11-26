package openmeteo_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"weather-server/internal/domain"
	openmeteo "weather-server/internal/forecast_providers/open_meteo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Response struct {
	Latitude  float64
	Longitude float64
	Timezone  string
	Daily     struct {
		Time        []string
		Temperature []float64 `json:"temperature_2m_max"`
	}
}

func TestGetForecastSuccess(t *testing.T) {
	t.Skip("Enable to test the real client")
	client := openmeteo.NewClient()

	ctx := context.Background()
	now := time.Now()
	q := domain.DayForecastQuery{
		Location: domain.Location{
			Latitude:  -13.52264,
			Longitude: -71.96734,
		},
		Day: time.Now(),
	}

	forecastRaw, err := client.GetDayForecast(ctx, q)
	require.NoError(t, err)
	marshalled, err := json.Marshal(forecastRaw)
	require.NoError(t, err)

	var resp Response
	require.NoError(t, json.Unmarshal(marshalled, &resp))
	assert.Equal(t, now.Format("2006-01-02"), resp.Daily.Time[0])
	assert.Equal(t, "GMT", resp.Timezone)
	t.Logf("%+v\n", resp)
}

func TestClientImplements(t *testing.T) {
	require.Implements(t, (*domain.Client)(nil), new(openmeteo.Client))
}

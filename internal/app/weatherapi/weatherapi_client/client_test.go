package weatherapiclient_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"weather-server/internal/app/weatherapi"
	weatherapiclient "weather-server/internal/app/weatherapi/weatherapi_client"
	"weather-server/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Change to a valid value for testing
const apiKey = "XD"

type Response struct {
	Location struct {
		Name    string
		Country string
		Lat     float64
		Lon     float64
	}
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	}
	Forecast struct {
		ForecastDay []struct {
			Date string
		}
	}
}

func TestGetForecastSuccess(t *testing.T) {
	t.Skip("Enable to test the real client")
	client := weatherapiclient.NewClient(apiKey)

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

	var resp Response
	require.NoError(t, json.Unmarshal(forecastRaw, &resp))
	assert.Equal(t, "Cusco", resp.Location.Name)
	assert.Equal(t, "Peru", resp.Location.Country)
	assert.Equal(t, now.Format("2006-01-02"), resp.Forecast.ForecastDay[0].Date)
	assert.Greater(t, resp.Current.TempF, resp.Current.TempC)
	t.Logf("%+v\n", resp)
}

func TestClientImplements(t *testing.T) {
	require.Implements(t, (*weatherapi.Client)(nil), new(weatherapiclient.Client))
}

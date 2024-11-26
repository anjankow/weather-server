package services_test

import (
	"context"
	"testing"
	"time"
	"weather-server/internal/domain"
	"weather-server/internal/forecast_providers/mock"
	openmeteo "weather-server/internal/forecast_providers/open_meteo"
	"weather-server/internal/forecast_providers/weatherapi"
	"weather-server/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAggregator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock1 := mock.NewMockClient(mockCtrl)
	mock2 := mock.NewMockClient(mockCtrl)

	provider1 := "mock 1"
	provider2 := "mock 2"

	ctx := context.Background()
	numOfDays := 33

	aggr, err := services.NewAggregator(
		numOfDays,
		map[string]domain.Client{
			provider1: mock1,
			provider2: mock2,
		})
	require.NoError(t, err)

	q := domain.Location{
		Latitude:  -13.52264,
		Longitude: -71.96734,
	}

	expected := map[string]domain.DayForecastRaw{
		provider1: {"data": 55},
		provider2: {"data": 33},
	}
	mock1.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return(expected[provider1], nil)
	mock2.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return(expected[provider2], nil)

	forecast, err := aggr.GetForecast(ctx, q)
	require.NoError(t, err)

	assert.Len(t, forecast, 2)
	checkForecastData := func(idx int) {
		expectedForecastData := expected[forecast[idx].APIName]
		assert.Equal(t, expectedForecastData, forecast[idx].DayForecasts[0])
	}

	checkForecastData(0)
	checkForecastData(1)
}

func TestAggregatorReal(t *testing.T) {
	t.Skip("Enable to test the real clients")

	weatherAPIKey := "fill me"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	numOfDays := 15

	aggr, err := services.NewAggregator(
		numOfDays,
		map[string]domain.Client{
			weatherapi.ProviderName: weatherapi.NewClient(weatherAPIKey),
			openmeteo.ProviderName:  openmeteo.NewClient(),
		})
	require.NoError(t, err)

	q := domain.Location{
		Latitude:  -13.52264,
		Longitude: -71.96734,
	}

	forecast, err := aggr.GetForecast(ctx, q)
	require.NoError(t, err)

	assert.Len(t, forecast, 2)
}

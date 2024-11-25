package app_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"weather-server/internal/app"
	"weather-server/internal/app/mock"
	"weather-server/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAggregator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock1 := mock.NewMockClient(mockCtrl)
	service1 := app.NewService(mock1)

	mock2 := mock.NewMockClient(mockCtrl)
	service2 := app.NewService(mock2)

	provider1 := "mock 1"
	provider2 := "mock 2"
	aggr := app.NewAggregator().AddService(provider1, service1).AddService(provider2, service2)

	ctx := context.Background()
	numOfDays := 33
	q := domain.ForecastQuery{
		Location: domain.Location{
			Latitude:  -13.52264,
			Longitude: -71.96734,
		},
		FromDay:   time.Now(),
		NumOfDays: numOfDays,
	}

	expected := map[string]domain.DayForecastRaw{
		provider1: domain.DayForecastRaw(json.RawMessage([]byte{55})),
		provider2: domain.DayForecastRaw(json.RawMessage([]byte{33})),
	}
	mock1.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return([]byte(expected[provider1]), nil)
	mock2.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return([]byte(expected[provider2]), nil)

	forecast, err := aggr.GetForecast(ctx, q)
	require.NoError(t, err)

	assert.Len(t, forecast, 2)
	checkForecastData := func(idx int) {
		expectedForecastData := expected[forecast[idx].APIName]
		assert.Equal(t, expectedForecastData[0], forecast[idx].DayForecasts[0][0])
	}

	checkForecastData(0)
	checkForecastData(1)
}

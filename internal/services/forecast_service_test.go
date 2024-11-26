package services_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"weather-server/internal/domain"
	"weather-server/internal/forecast_providers/mock"
	"weather-server/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetForecastSuccessMultiple(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockClient(mockCtrl)
	service := services.NewForecastService(mock)

	ctx := context.Background()
	numOfDays := 450
	q := domain.ForecastQuery{
		Location: domain.Location{
			Latitude:  -13.52264,
			Longitude: -71.96734,
		},
		FromDay:   time.Now(),
		NumOfDays: numOfDays,
	}

	mock.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return(domain.DayForecastRaw{"data": 55}, nil)

	forecast, err := service.GetForecast(ctx, q)
	require.NoError(t, err)
	require.Len(t, forecast, numOfDays)

	assert.Equal(t, forecast[numOfDays-1], domain.DayForecastRaw{
		"data": 55,
	})
}

func TestGetForecastSuccessOne(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockClient(mockCtrl)
	service := services.NewForecastService(mock)

	ctx := context.Background()
	q := domain.ForecastQuery{
		Location: domain.Location{
			Latitude:  -13.52264,
			Longitude: -71.96734,
		},
		FromDay:   time.Now(),
		NumOfDays: 1,
	}

	mock.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(1).Return(domain.DayForecastRaw{"data": 55}, nil)

	forecast, err := service.GetForecast(ctx, q)
	require.NoError(t, err)
	require.Len(t, forecast, 1)

	assert.Equal(t, forecast[0], domain.DayForecastRaw{"data": 55})
}

func TestGetForecastClientFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockClient(mockCtrl)
	service := services.NewForecastService(mock)

	ctx := context.Background()
	q := domain.ForecastQuery{
		Location: domain.Location{
			Latitude:  -13.52264,
			Longitude: -71.96734,
		},
		FromDay:   time.Now(),
		NumOfDays: 1,
	}

	testError := errors.New("client error")
	mock.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).
		Return(nil, testError)

	_, err := service.GetForecast(ctx, q)
	require.ErrorAs(t, err, &testError)
}

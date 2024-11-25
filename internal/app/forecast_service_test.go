package app_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
	"weather-server/internal/app"
	"weather-server/internal/app/mock"
	"weather-server/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetForecastSuccessMultiple(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockClient(mockCtrl)
	service := app.NewService(mock)

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

	mock.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return([]byte{55}, nil)

	forecast, err := service.GetForecast(ctx, q)
	require.NoError(t, err)
	require.Len(t, forecast, numOfDays)

	assert.Equal(t, forecast[numOfDays-1], domain.DayForecastRaw(json.RawMessage([]byte{55})))
}

func TestGetForecastSuccessOne(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockClient(mockCtrl)
	service := app.NewService(mock)

	ctx := context.Background()
	q := domain.ForecastQuery{
		Location: domain.Location{
			Latitude:  -13.52264,
			Longitude: -71.96734,
		},
		FromDay:   time.Now(),
		NumOfDays: 1,
	}

	mock.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(1).Return([]byte{55}, nil)

	forecast, err := service.GetForecast(ctx, q)
	require.NoError(t, err)
	require.Len(t, forecast, 1)

	assert.Equal(t, forecast[0], domain.DayForecastRaw(json.RawMessage([]byte{55})))
}

func TestGetForecastClientFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockClient(mockCtrl)
	service := app.NewService(mock)

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

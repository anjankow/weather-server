package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weather-server/internal/api/handlers"
	"weather-server/internal/api/server"
	"weather-server/internal/app"
	"weather-server/internal/domain"
	"weather-server/internal/forecast_providers/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetWeatherOneProvider(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mock := mock.NewMockClient(mockCtrl)
	
	numOfDays := 5
	mock.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return([]byte{55}, nil)
	s := NewTestServer(t, numOfDays, map[string]domain.Client{
		"mock1": mock,
	})

	res := MakeRequest(t, s, "/weather")

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	var resJson json.RawMessage
	require.NoError(t,resJson.UnmarshalJSON(resJson))

	t.Log(string(resJson))
}


func MakeRequest(t *testing.T, s server.Server, url string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/-/healthy", nil)
	res := httptest.NewRecorder()

	s.Echo.ServeHTTP(res, req)
	return res
}

func NewTestServer(t *testing.T, numOfDays int, providers map[string]domain.Client) server.Server{
	t.Helper()

	app, err := app.New(app.Config{
		NumOfForecastDays: numOfDays,
	}, app.Dependencies{
		ForecastProviders: providers,
	})
	require.NoError(t, err)

	s := server.New(server.Config{
		Debug: true,
		RequestTimeout: 10*time.Second,
	})
	handlers.AttachHandlers(s, app)
	
	return s
}
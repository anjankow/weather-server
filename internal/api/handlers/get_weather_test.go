package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	numOfDays := 4
	type forecastType struct {
		Temp    float64
		Comment string
	}
	expectedForecast := forecastType{
		Temp:    -22.6,
		Comment: "never go outside",
	}

	expectedForecastJson, err := json.Marshal(expectedForecast)
	require.NoError(t, err)
	var expectedForecastMap domain.DayForecastRaw
	require.NoError(t, json.Unmarshal(expectedForecastJson, &expectedForecastMap))

	// make sure that GetDayForecast gets called for each day
	mock.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return(expectedForecastMap, nil)

	// add just one forecast provider
	s := NewTestServer(t, numOfDays, map[string]domain.Client{
		"mock1": mock,
	})

	res := MakeRequest(t, s, "/weather?longitude=33&latitude=22")

	require.Equal(t, http.StatusOK, res.Result().StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	t.Log(string(body))

	// make sure that there's just one forecast provider in the response
	var resJson map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &resJson))
	require.Len(t, resJson, 1)

	// assert the provider name
	forecast, ok := resJson["mock1"]
	require.True(t, ok)

	// make sure that there are numOfDays days listed with the correct numbers
	var forecastJson map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(forecast, &forecastJson))
	require.Len(t, forecastJson, numOfDays)
	day1, ok := forecastJson["day1"]
	require.True(t, ok)
	var forecast1 forecastType
	require.NoError(t, json.Unmarshal(day1, &forecast1))
	// check the forecast content
	assert.Equal(t, expectedForecast, forecast1)

	day4, ok := forecastJson["day4"]
	require.True(t, ok)
	var forecast4 forecastType
	require.NoError(t, json.Unmarshal(day4, &forecast4))
	// check the forecast content
	assert.Equal(t, expectedForecast, forecast4)

}

func TestGetWeatherMultipleProviders(t *testing.T) {
	type forecastType struct {
		Temp    float64
		Comment string
	}
	numOfDays := 4

	setMockExpectations := func(forecast forecastType, m *mock.MockClient) {
		expectedForecastJson, err := json.Marshal(forecast)
		require.NoError(t, err)
		var expectedForecastMap domain.DayForecastRaw
		require.NoError(t, json.Unmarshal(expectedForecastJson, &expectedForecastMap))

		m.EXPECT().GetDayForecast(gomock.Any(), gomock.Any()).Times(numOfDays).Return(expectedForecastMap, nil)
	}

	// create 3 different forecast providers
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mock1 := mock.NewMockClient(mockCtrl)
	mock2 := mock.NewMockClient(mockCtrl)
	mock3 := mock.NewMockClient(mockCtrl)

	setMockExpectations(forecastType{
		Temp:    -22.6,
		Comment: "never go outside",
	}, mock1)

	setMockExpectations(forecastType{
		Temp:    13.0,
		Comment: "it's alright",
	}, mock2)

	setMockExpectations(forecastType{
		Temp:    44,
		Comment: "perfect",
	}, mock3)

	// create a server with 3 mocked providers
	s := NewTestServer(t, numOfDays, map[string]domain.Client{
		"mock1": mock1,
		"mock2": mock2,
		"mock3": mock3,
	})

	res := MakeRequest(t, s, "/weather?longitude=33&latitude=22")

	require.Equal(t, http.StatusOK, res.Result().StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	t.Log(string(body))

	var resJson map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &resJson))
	require.Len(t, resJson, 3)

	// make sure that all the providers are included in the response
	_, ok := resJson["mock1"]
	require.True(t, ok)
	forecast2, ok := resJson["mock2"]
	require.True(t, ok)
	forecast3, ok := resJson["mock3"]
	require.True(t, ok)

	// make sure that the responses of each provider differ
	{
		var forecastJson map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(forecast2, &forecastJson))
		require.Len(t, forecastJson, numOfDays)
		day1, ok := forecastJson["day3"]
		require.True(t, ok)
		var forecast forecastType
		require.NoError(t, json.Unmarshal(day1, &forecast))
		assert.Equal(t, "it's alright", forecast.Comment)
	}

	{
		var forecastJson map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(forecast3, &forecastJson))
		require.Len(t, forecastJson, numOfDays)
		day1, ok := forecastJson["day3"]
		require.True(t, ok)
		var forecast forecastType
		require.NoError(t, json.Unmarshal(day1, &forecast))
		assert.Equal(t, "perfect", forecast.Comment)
	}
}

func TestGetWeatherMissingParams(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mock := mock.NewMockClient(mockCtrl)

	numOfDays := 4

	s := NewTestServer(t, numOfDays, map[string]domain.Client{
		"mock1": mock,
	})

	res := MakeRequest(t, s, "/weather")

	assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	t.Log(string(body))
}

func TestGetWeatherInvalidParams(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mock := mock.NewMockClient(mockCtrl)

	numOfDays := 4

	s := NewTestServer(t, numOfDays, map[string]domain.Client{
		"mock1": mock,
	})

	res := MakeRequest(t, s, "/weather?longitude=181&latitude=77")

	assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	t.Log(string(body))
}

func MakeRequest(t *testing.T, s server.Server, url string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, url, nil)
	res := httptest.NewRecorder()

	s.Echo.ServeHTTP(res, req)
	return res
}

func NewTestServer(t *testing.T, numOfDays int, providers map[string]domain.Client) server.Server {
	t.Helper()

	app, err := app.New(app.Config{
		NumOfForecastDays: numOfDays,
	}, app.Dependencies{
		ForecastProviders: providers,
	})
	require.NoError(t, err)

	s := server.New(server.Config{
		Debug:          true,
		RequestTimeout: 10 * time.Second,
	}, app)

	return s
}

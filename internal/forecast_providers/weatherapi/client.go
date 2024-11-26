package weatherapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"weather-server/internal/domain"
)

const ProviderName = "WeatherAPI"

const baseURL = "https://api.weatherapi.com/v1/"

type Client struct {
	apiKey string
}

func NewClient(apiKey string) Client {
	return Client{apiKey}
}

func (c Client) GetDayForecast(ctx context.Context, query domain.DayForecastQuery) (domain.DayForecastRaw, error) {
	url, err := url.JoinPath(baseURL, "forecast.json")
	if err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to parse weatherapi forecast url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.DayForecastRaw{}, err
	}
	// add query params
	q := req.URL.Query()
	q.Set("key", c.apiKey)
	// Latitude and Longitude (Decimal degree) e.g: q=48.8567,2.3508
	location := fmt.Sprintf("%.4f,%.4f", query.Latitude, query.Longitude)
	q.Set("q", location)
	q.Set("date", query.Day.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	// do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to exec weatherapi req: %w", err)
	}
	defer resp.Body.Close()

	// read the response, in case of failure it contains the error details
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to read weatherapi resp: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden {
			return domain.DayForecastRaw{}, errors.New("invalid weatherapi API key")
		}

		type weatherapiError struct {
			Error struct {
				Message string
			}
		}

		var respError weatherapiError
		if err := json.Unmarshal(body, &respError); err != nil {
			return domain.DayForecastRaw{}, fmt.Errorf("weatherapi failed with status: %s", resp.Status)
		}

		return domain.DayForecastRaw{}, fmt.Errorf("weatherapi request failed: %s", respError.Error.Message)
	}

	var content domain.DayForecastRaw
	if err := json.Unmarshal(body, &content); err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return content, nil
}

package weatherapiclient

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

const baseURL = "https://api.weatherapi.com/v1/"

type Client struct {
	apiKey string
}

func NewClient(apiKey string) Client {
	return Client{apiKey}
}

func (c Client) GetDayForecast(ctx context.Context, query domain.DayForecastQuery) (json.RawMessage, error) {
	url, err := url.JoinPath(baseURL, "forecast.json")
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("failed to parse weatherapi forecast url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return json.RawMessage{}, err
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
		return json.RawMessage{}, fmt.Errorf("failed to exec weatherapi req: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden {
			return json.RawMessage{}, errors.New("invalid weatherapi API key")
		}
		return json.RawMessage{}, fmt.Errorf("weatherapi failed with status: %s", resp.Status)
	}

	// read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return json.RawMessage{}, fmt.Errorf("failed to read weatherapi resp: %w", err)
	}

	return body, nil
}

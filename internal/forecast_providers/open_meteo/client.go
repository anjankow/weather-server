package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"weather-server/internal/domain"
)

const ProviderName = "OpenMeteo"

const baseURL = "https://api.open-meteo.com/v1/"

type Client struct {
}

func NewClient() Client {
	return Client{}
}

// GetDayForecast gets up to 15 days of forecast
func (c Client) GetDayForecast(ctx context.Context, query domain.DayForecastQuery) (domain.DayForecastRaw, error) {
	url, err := url.JoinPath(baseURL, "forecast")
	if err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to parse openmeteo forecast url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.DayForecastRaw{}, err
	}
	// add query params
	q := req.URL.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", query.Latitude))
	q.Set("longitude", fmt.Sprintf("%.4f", query.Longitude))
	date := query.Day.Format("2006-01-02")
	q.Set("start_date", date)
	q.Set("end_date", date)
	q.Set("daily", "temperature_2m_max")
	req.URL.RawQuery = q.Encode()

	// do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to exec openmeteo req: %w", err)
	}
	defer resp.Body.Close()

	// read the response, in case of failure it contains the error details
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to read openmeteo resp: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		type openmeteoError struct {
			Reason string
		}

		var respError openmeteoError
		if err := json.Unmarshal(body, &respError); err != nil {
			return domain.DayForecastRaw{}, fmt.Errorf("openmeteo failed with status: %s", resp.Status)
		}

		return domain.DayForecastRaw{}, fmt.Errorf("openmeteo request failed: %s", respError.Reason)
	}

	var content domain.DayForecastRaw
	if err := json.Unmarshal(body, &content); err != nil {
		return domain.DayForecastRaw{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return content, nil
}

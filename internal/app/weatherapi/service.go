package weatherapi

import (
	"context"
	"errors"
	"fmt"
	"time"
	"weather-server/internal/domain"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	client Client
}

func NewService(client Client) Service {
	return Service{
		client: client,
	}
}

func (s Service) GetForecast(ctx context.Context, query domain.ForecastQuery) (domain.ForecastResponse, error) {

	type dayForecast struct {
		DayNum int
		domain.DayForecastRaw
	}
	numOfDays := query.NumOfDays

	forecastChan := make(chan dayForecast, numOfDays)
	g, gctx := errgroup.WithContext(ctx)

	// Get forecast for the numOfDays consequitive days in parallel
	for dayNum := 1; dayNum <= numOfDays; dayNum++ {
		day := query.FromDay.Add(time.Hour * 24 * time.Duration(dayNum))
		dayQuery := domain.DayForecastQuery{
			Location: query.Location,
			Day:      day,
		}

		g.Go(func() error {
			resp, err := s.client.GetDayForecast(gctx, dayQuery)
			if err != nil {
				return err
			}
			if len(resp) == 0 {
				return errors.New("failed to get weatherapi day forecast: empty body")
			}

			// No failures, add the response to the channel
			forecastChan <- dayForecast{
				DayNum:         dayNum,
				DayForecastRaw: domain.DayForecastRaw(resp),
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("weatherapi service failed: GetForecast: %w", err)
	}

	ret := make(domain.ForecastResponse, query.NumOfDays)
	for f := range forecastChan {
		// Return the daily forecasts in order
		// Index 0 is the 1st day
		ret[f.DayNum-1] = f.DayForecastRaw
	}

	return ret, nil
}

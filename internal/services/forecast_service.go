package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"weather-server/internal/domain"

	"golang.org/x/sync/errgroup"
)

type ForecastService struct {
	client domain.Client
}

func NewForecastService(client domain.Client) ForecastService {
	return ForecastService{
		client: client,
	}
}

func (s ForecastService) GetForecast(ctx context.Context, query domain.ForecastQuery) (domain.DayForecastSlice, error) {
	if query.NumOfDays <= 0 {
		return nil, errors.New("invalid number of days")
	}

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
				return errors.New("failed to get a day forecast: empty body")
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
		return nil, err
	}

	if len(forecastChan) != query.NumOfDays {
		return nil, fmt.Errorf("number of received forecasts (%v) is not as requested (%v)", len(forecastChan), query.NumOfDays)
	}

	dayForecasts := make(domain.DayForecastSlice, query.NumOfDays)
	for i := 0; i < query.NumOfDays; i++ {
		// Return the daily forecasts in order
		// Index 0 is the 1st day
		var f = <-forecastChan
		dayForecasts[f.DayNum-1] = f.DayForecastRaw
	}

	return dayForecasts, nil
}

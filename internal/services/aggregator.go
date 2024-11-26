package services

import (
	"context"
	"errors"
	"fmt"
	"weather-server/internal/domain"

	"golang.org/x/sync/errgroup"
)

type Aggregator struct {
	forecastServices map[string]ForecastService
}

func NewAggregator(forecastProviders map[string /*provider name*/]domain.Client) (Aggregator,error) {
	if forecastProviders == nil {
		return Aggregator{}, errors.New("empty providers set")
	}

	services := make(map[string]ForecastService, len(forecastProviders))
	for providerName, provider := range forecastProviders {
		services[providerName] = NewForecastService(provider)
	}
	
	return Aggregator{
		forecastServices: services,
	}, nil
}


func (a Aggregator) GetForecast(ctx context.Context, query domain.ForecastQuery) (domain.ForecastAggregate, error) {

	forecastsChan := make(chan domain.ForecastResponse, len(a.forecastServices))
	g, gctx := errgroup.WithContext(ctx)

	for provider, service := range a.forecastServices {
		g.Go(func() error {
			forecast, err := service.GetForecast(gctx, query)
			if err != nil {
				return fmt.Errorf("%s provider failed: %w", provider, err)
			}
			forecastsChan <- domain.ForecastResponse{
				DayForecasts: forecast,
				APIName:      provider,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if len(forecastsChan) != len(a.forecastServices) {
		return nil, fmt.Errorf("number of received forecasts (%v) is not equal to number of forecast services (%v)", len(forecastsChan), len(a.forecastServices))
	}

	aggregate := make(domain.ForecastAggregate, 0, len(a.forecastServices))
	for i := 0; i < len(a.forecastServices); i++ {
		aggregate = append(aggregate, <-forecastsChan)
	}

	return aggregate, nil
}

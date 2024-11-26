package main

import (
	"fmt"
	"os"
	"weather-server/internal/api/server"
	"weather-server/internal/app"
	"weather-server/internal/domain"
	openmeteo "weather-server/internal/forecast_providers/open_meteo"
	"weather-server/internal/forecast_providers/weatherapi"

	"github.com/joho/godotenv"
)

func main() {
	appCfg := app.Config{
		NumOfForecastDays: 5,
	}

	// read the WeatherAPI key from env
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to read the .env file: %s", err.Error())
		os.Exit(1)
	}
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		fmt.Printf("Set the WEATHER_API_KEY in .evn file")
		os.Exit(1)
	}

	// initialize app dependencies
	appDeps := app.Dependencies{
		ForecastProviders: map[string]domain.Client{
			weatherapi.ProviderName: weatherapi.NewClient(weatherAPIKey),
			openmeteo.ProviderName:  openmeteo.NewClient(),
		},
	}

	// initialize the application
	app, err := app.New(appCfg, appDeps)
	if err != nil {
		fmt.Printf("Failed to initialize App: %s", err.Error())
		os.Exit(1)
	}

	// run the server
	server := server.New(
		server.Config{
			Debug: false,
			ListenAddr: ":8099",
		},
		app,
	)
	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start server: %s", err.Error())
		os.Exit(1)
	}
}

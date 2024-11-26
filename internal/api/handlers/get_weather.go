package handlers

import (
	"fmt"
	"net/http"
	"weather-server/internal/api/httperrors"
	"weather-server/internal/app"
	"weather-server/internal/domain"

	"github.com/labstack/echo/v4"
)

func GetWeather(app app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// get and validate the query params
		var params = struct {
			Latitude  float64
			Longitude float64
		}{}
		binder := c.Echo().Binder.(*echo.DefaultBinder)
		if err := binder.BindQueryParams(c, &params); err != nil {
			return httperrors.NewValidationError(err.Error())
		}

		if params.Longitude > domain.LongitudeMax || params.Longitude < domain.LongitudeMin {
			return httperrors.NewValidationError("invalid longitude value")
		}

		if params.Latitude > domain.LatitudeMax || params.Latitude < domain.LatitudeMin {
			return httperrors.NewValidationError("invalid latitude value")
		}

		// map the params to the domain objects
		location := domain.Location{
			Longitude: params.Longitude,
			Latitude:  params.Latitude,
		}
		// call the respective app service
		forecast, err := app.Aggregator.GetForecast(ctx, location)
		if err != nil {
			return httperrors.NewInternalServerError(err)
		}

		// map the domain object to response api
		//! {
		//     "weatherAPI1": {
		//       "day1": {...},
		//       "day2": {...},
		//       ...
		//     },
		//     "weatherAPI2": {
		//       "day1": {...},
		//       "day2": {...},
		//       ...
		//     }
		//!   }
		resp := make(map[string]interface{}, len(forecast))
		for _, f := range forecast {
			dayForecastResp := make(map[string]interface{}, len(f.DayForecasts))
			for i, df := range f.DayForecasts {
				// index 0 in DayForecasts maps to day1
				day := fmt.Sprintf("day%v", i+1)
				dayForecastResp[day] = df
			}
			resp[f.APIName] = dayForecastResp
		}

		return c.JSON(http.StatusOK, resp)
	}
}

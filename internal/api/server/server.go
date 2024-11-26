package server

import (
	"errors"
	"log"
	"net/http"
	"time"
	"weather-server/internal/api/handlers"
	"weather-server/internal/api/httperrors"
	"weather-server/internal/app"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Echo *echo.Echo
	Cfg  Config
}

type Config struct {
	Debug          bool
	RequestTimeout time.Duration
	MgmtKey        string
}

func New(cfg Config, app app.App) Server {
	s := Server{
		Echo: newRouter(cfg, app),
		Cfg:  cfg,
	}

	return s
}

func newRouter(cfg Config, app app.App) *echo.Echo {
	e := echo.New()
	e.Debug = cfg.Debug
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	if cfg.RequestTimeout > 0 {
		e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: cfg.RequestTimeout,
		}))
	}
	e.HTTPErrorHandler = func(err error, c echo.Context) {

		var httpError httperrors.HTTPError

		if errors.As(err, &httpError) && httpError.StatusCode >= http.StatusInternalServerError {
			log.Default().Printf("internal server error: %s\n", err.Error())
			// hide the error details
			httpError.StatusCode = http.StatusInternalServerError
			httpError.Err = errors.New("something went wrong, our monkeys are already investigating")

			if !c.Response().Committed {
				if c.Request().Method == http.MethodHead {
					err = c.NoContent(httpError.StatusCode)
				} else {
					err = c.String(httpError.StatusCode, httpError.Err.Error())
				}

				if err != nil {
					log.Default().Printf("failed to handle http error: %s\n", err.Error())
				}
			}
			return
		}

		e.DefaultHTTPErrorHandler(err, c)
	}

	e.GET("/weather", handlers.GetWeather(app))

	mgmtGroup := e.Group("/-")
	if cfg.MgmtKey != "" {
		mgmtGroup.Use(middleware.KeyAuth(func(auth string, _ echo.Context) (bool, error) {
			return auth == cfg.MgmtKey, nil
		}))
	}
	mgmtGroup.GET("/healthy", handlers.GetHealthy())

	return e
}

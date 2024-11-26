package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Echo *echo.Echo
	Cfg Config
}

type Config struct {
	Debug          bool
	RequestTimeout time.Duration
	MgmtKey        string
}


func New(cfg Config) Server {
	s := Server{
		Echo: newRouter(cfg),
	}
	
	return s
}

func newRouter(cfg Config) *echo.Echo{
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

	return e
}
package handlers

import (
	"weather-server/internal/api/server"
	"weather-server/internal/app"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AttachHandlers(s server.Server, app app.App) {
	e := s.Echo

	root := e.Group("")
	root.GET("/weather", GetWeather(app))

	mgmtGroup := e.Group("/-")
	if s.Cfg.MgmtKey != "" {
		mgmtGroup.Use(middleware.KeyAuth(func(auth string, _ echo.Context) (bool, error) {
			return auth == s.Cfg.MgmtKey, nil
		}))
	}
	mgmtGroup.GET("/healthy", GetHealthy())

}
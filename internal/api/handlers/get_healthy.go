package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHealthy() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm doing great, thanks.")
	}
}

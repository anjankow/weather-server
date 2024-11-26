package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-server/internal/api/server"
	"weather-server/internal/app"

	"github.com/stretchr/testify/assert"
)

func TestGetHealthy(t *testing.T) {
	mgmtKey := "xxx"
	s := server.New(server.Config{Debug: true, MgmtKey: mgmtKey}, app.App{})

	req := httptest.NewRequest(http.MethodGet, "/-/healthy", nil)
	req.Header.Add("Authorization", "Bearer "+mgmtKey)

	res := httptest.NewRecorder()
	s.Echo.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
}

func TestGetHealthyUnauthorized(t *testing.T) {
	mgmtKey := "xxx"
	s := server.New(server.Config{Debug: true, MgmtKey: mgmtKey}, app.App{})

	req := httptest.NewRequest(http.MethodGet, "/-/healthy", nil)
	req.Header.Add("Authorization", "Bearer "+"invalid")

	res := httptest.NewRecorder()
	s.Echo.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
}

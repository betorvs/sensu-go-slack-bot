package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

// Health struct
type Health struct {
	Status string `json:"status"`
}

// CheckHealth func to return OK and http 200
func CheckHealth(c echo.Context) error {
	health := Health{}
	health.Status = "UP"
	return c.JSON(http.StatusOK, health)
}

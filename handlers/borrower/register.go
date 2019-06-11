package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func RegisterBorrower(c echo.Context) error {
	return c.JSON(http.StatusOK, "nOICE")
}

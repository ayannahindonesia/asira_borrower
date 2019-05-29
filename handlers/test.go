package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Test(c echo.Context) error {
	type TestResponse struct {
		Response string `json:"resp"`
	}

	testresponse := TestResponse{
		Response: "hello",
	}
	return c.JSON(http.StatusOK, testresponse)
}

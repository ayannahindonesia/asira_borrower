package handlers

import (
	"asira_borrower/asira"
	"net/http"

	"github.com/labstack/echo"
)

func AsiraAppInfo(c echo.Context) error {
	type AppInfo struct {
		AppName string                 `json:"app_name"`
		Version string                 `json:"version"`
		ENV     string                 `json:"env"`
		Config  map[string]interface{} `json:"configs"`
	}

	var show AppInfo

	show.AppName = asira.App.Name
	show.Version = asira.App.Version
	show.ENV = asira.App.ENV
	show.Config = asira.App.Config.GetStringMap(asira.App.ENV)

	return c.JSON(http.StatusOK, show)
}

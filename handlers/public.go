package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func OpenAndroidDeeplinkAsset(c echo.Context) error {
	gopath, _ := os.Getwd()
	jsonFile, err := os.Open(gopath + "/assets/assetlinks.json")

	defer jsonFile.Close()

	if err != nil {
		log.Println(err)

		return returnInvalidResponse(http.StatusInternalServerError, "", "error opening assetlinks.json")
	}

	var i interface{}
	byteV, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteV, &i)

	return c.JSON(http.StatusOK, i)
}

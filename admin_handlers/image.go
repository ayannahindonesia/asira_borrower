package admin_handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetImageB64String(c echo.Context) error {
	defer c.Request().Body.Close()
	image := models.Image{}

	imageId, _ := strconv.Atoi(c.Param("image_id"))
	_, err := image.FindbyID(imageId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Gambar Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, image)
}

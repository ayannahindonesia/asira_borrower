package admin_handlers

import (
	"asira_borrower/models"
	"net/http"

	guuid "github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func CreateClientConfig(c echo.Context) error {
	defer c.Request().Body.Close()

	clientC := models.Client_config{}

	payloadRules := govalidator.MapData{
		"name": []string{"required"},
		"key":  []string{"required"},
		"role": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &clientC)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	clientC.Secret = guuid.New().String()
	newConfig, err := clientC.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Client Config")
	}

	return c.JSON(http.StatusCreated, newConfig)
}

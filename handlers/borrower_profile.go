package handlers

import (
	"asira/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

func BorrowerProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	borrower, err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	return c.JSON(http.StatusOK, borrower)
}

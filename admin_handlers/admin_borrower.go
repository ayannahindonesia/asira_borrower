package admin_handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func BorrowerGetAll(c echo.Context) error {
	defer c.Request().Body.Close()

	borrower := models.Borrower{}
	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")
	// filters
	fullname := c.QueryParam("fullname")

	type Filter struct {
		Fullname string `json:"fullname" condition:"LIKE"`
	}
	result, err := borrower.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Fullname: fullname,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Borrower tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func BorrowerGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(c.Param("borrower_id"))
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Borrower ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}

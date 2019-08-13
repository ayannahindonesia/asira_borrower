package admin_handlers

import (
	"asira_borrower/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func LoanGetAll(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}
	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")
	owner := c.QueryParam("owner")

	type Filter struct {
		Owner sql.NullInt64 `json:"owner"`
	}
	result, err := loan.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Owner: sql.NullInt64{
			Int64: int64(owner),
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Loan tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func LoanGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	loanModel := models.Loan{}

	loanID, _ := strconv.Atoi(c.Param("loan_id"))
	loan, err := loanModel.FindbyID(loanID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Loan ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, loan)
}

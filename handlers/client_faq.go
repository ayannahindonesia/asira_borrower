package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//FAQPayload payload
type FAQPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

//FAQList get FAQ list
func FAQList(c echo.Context) error {
	defer c.Request().Body.Close()
	if err != nil {
		NLog("warning", "FAQList", fmt.Sprintf("unauthorized access FAQList : '%v'", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	var (
		faq    models.FAQ
		result basemodel.PagedFindResult
	)

	if searchAll := c.QueryParam("search_all"); len(searchAll) > 0 {
		type Filter struct {
			Title       string `json:"title" condition:"LIKE,optional"`
			Description string `json:"description" condition:"LIKE,optional"`
		}
		result, err = faq.PagedFindFilter(page, rows, orderby, sort, &Filter{
			Title:       searchAll,
			Description: searchAll,
		})
	} else {
		type Filter struct {
			Title       string `json:"title" condition:"LIKE"`
			Description string `json:"description"`
		}
		result, err = faq.PagedFindFilter(page, rows, orderby, sort, &Filter{
			Title:       c.QueryParam("title"),
			Description: c.QueryParam("description"),
		})
	}

	if err != nil {
		NLog("warning", "FAQList", fmt.Sprintf("error finding FAQ : %v", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

// FAQDetail get FAQ detail by id
func FAQDetail(c echo.Context) error {
	defer c.Request().Body.Close()
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	faqID, _ := strconv.ParseUint(c.Param("faq_id"), 10, 64)

	faq := models.FAQ{}
	err = faq.FindbyID(faqID)
	if err != nil {
		NLog("warning", "FAQDetail", fmt.Sprintf("FAQ %v not found : %v", faqID, err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusNotFound, err, "Tidak memiliki hak akses")
	}

	return c.JSON(http.StatusOK, faq)
}

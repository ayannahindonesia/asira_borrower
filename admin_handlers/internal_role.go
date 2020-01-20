package admin_handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func AddRole(c echo.Context) error {
	defer c.Request().Body.Close()

	Iroles := models.InternalRoles{}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"system":      []string{"required"},
		"status":      []string{},
		"description": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &Iroles)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err := Iroles.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Internal Roles")
	}

	return c.JSON(http.StatusCreated, Iroles)
}

func UpdateRole(c echo.Context) error {
	defer c.Request().Body.Close()
	Iroles_id, _ := strconv.ParseUint(c.Param("role_id"), 10, 64)

	Iroles := models.InternalRoles{}
	err := Iroles.FindbyID(Iroles_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Internal Role %v tidak ditemukan", Iroles_id))
	}

	payloadRules := govalidator.MapData{
		"name":        []string{"required"},
		"system":      []string{"required"},
		"status":      []string{},
		"description": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &Iroles)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = Iroles.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update Internal Roles %v", Iroles_id))
	}

	return c.JSON(http.StatusOK, Iroles)
}

func GetAllRole(c echo.Context) error {
	defer c.Request().Body.Close()

	Iroles := models.InternalRoles{}
	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	var Filter struct{}
	result, err := Iroles.PagedFilterSearch(page, rows, orderby, sort, &Filter)

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Internal Role tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func RoleGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	Iroles := models.InternalRoles{}

	IrolesID, _ := strconv.ParseUint(c.Param("role_id"), 10, 64)
	err := Iroles.FindbyID(IrolesID)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Role ID tidak ditemukan")
	}

	return c.JSON(http.StatusOK, Iroles)
}

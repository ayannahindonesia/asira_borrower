package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

type (
	AgentPayload struct {
		Email string  `json:"email"`
		Phone string  `json:"phone"`
		Banks []int64 `json:"banks"`
		Image string  `json:"image"`
	}

	BanksResponse struct {
		Name string `json:"name"`
	}

	AgentResponse struct {
		models.Agent
		BankNames pq.StringArray `json:"bank_names"`
	}
)

//AgentProfile get current agent's profile
func AgentProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	agentModel := models.Agent{}

	agentID, _ := strconv.Atoi(claims["jti"].(string))
	err := agentModel.FindbyID(agentID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	return c.JSON(http.StatusOK, agentModel)
}

//AgentProfileEdit update current agent's profile
func AgentProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()
	var agentPayload AgentPayload

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.Atoi(claims["jti"].(string))

	//cek agent with custom field (name of banks)
	agentModel := models.Agent{}
	err := agentModel.FindbyID(agentID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	//securing old password
	password := agentModel.Password

	payloadRules := govalidator.MapData{
		"email": []string{"email"},
		"phone": []string{"id_phonenumber"},
		"banks": []string{"valid_id:banks"},
		"image": []string{},
	}

	//validate request data
	validate := validateRequestPayload(c, payloadRules, &agentPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek unique for patching
	uniques := map[string]string{
		"email": agentPayload.Email,
		"phone": agentPayload.Phone,
	}
	foundFields, err := checkPatchFields("agents", "id", agentModel.ID, uniques)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error : "+foundFields)
	}

	if len(agentPayload.Email) > 0 {
		agentModel.Email = agentPayload.Email
	}

	if len(agentPayload.Phone) > 0 {
		agentModel.Phone = agentPayload.Phone
	}

	//if payload not 0 and category must "agent" not "account_executive"
	if len(agentPayload.Banks) > 0 && agentModel.Category != "account_executive" {
		agentModel.Banks = pq.Int64Array(agentPayload.Banks)
	}

	if len(agentPayload.Image) > 0 {
		Image := models.Image{}
		//search if exist
		err = Image.FindbyID(int(agentModel.ImageID.Int64))
		Image.Image_string = agentPayload.Image

		if err != nil {
			err = Image.Create()
			if err != nil {
				return returnInvalidResponse(http.StatusInternalServerError, err, "Failed storing image")
			}
		} else {
			err = Image.Save()
			if err != nil {
				return returnInvalidResponse(http.StatusInternalServerError, err, "Failed storing image")
			}
		}

		//update id image
		agentModel.ImageID = sql.NullInt64{
			Int64: int64(Image.ID),
			Valid: true,
		}
	}
	//restoring old password and update data
	agentModel.Password = password
	err = agentModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Gagal mengubah data akun agen")
	}

	//Refetching after update
	var response AgentResponse
	db := asira.App.DB.Table("agents ag").
		Select("ag.*, (SELECT ARRAY_AGG(name) FROM banks WHERE id IN (SELECT UNNEST(ag.banks))) as bank_names").
		Where("ag.id = ?", agentID)
	err = db.Find(&response).Error
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	return c.JSON(http.StatusOK, response)
}

func checkPatchFields(tableName string, fieldID string, id uint64, uniques map[string]string) (string, error) {
	var count int
	fieldsFound := ""

	//...check unique
	for key, val := range uniques {
		//init query
		db := asira.App.DB
		db = db.Table(tableName).Select(fieldID)

		//get users other than idcardNumber...
		db = db.Not(fieldID, id)

		//if field not empty
		if len(val) > 0 || val != "" {
			db = db.Where(fmt.Sprintf("LOWER(%s) = ?", key), strings.ToLower(val))
		} else {
			//skip checking
			continue
		}
		//query count
		err = db.Count(&count).Error
		if err != nil || count > 0 {
			fieldsFound += key + ", "
		}
	}
	if fieldsFound != "" {
		return fieldsFound, errors.New("data unique already exist")
	}
	return fieldsFound, nil
}

package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

type AgentPayload struct {
	Email string  `json:"email"`
	Phone string  `json:"phone"`
	Banks []int64 `json:"banks"`
	Image string  `json:"image"`
}

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

		unbased, _ := base64.StdEncoding.DecodeString(agentPayload.Image)
		filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
		url, err := asira.App.S3.UploadJPEG(unbased, filename)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal upload foto agent")
		}

		//DONE: delete old image
		if len(agentModel.Image) > 0 {
			i := strings.Split(agentModel.Image, "/")
			delImage := i[len(i)-1]
			err = asira.App.S3.DeleteObject(delImage)
			if err != nil {
				log.Printf("failed to delete image %v from s3 bucket", delImage)
			}
		}

		agentModel.Image = url
	}
	//restoring old password and update data
	agentModel.Password = password
	err = agentModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Gagal mengubah data akun agen")
	}

	return c.JSON(http.StatusOK, agentModel)
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

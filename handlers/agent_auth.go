package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	//AgentLoginCreds payload agent's login credentials
	AgentLoginCreds struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
)

//AgentLogin borrower login, borrower can choose either login with email / phone
func AgentLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentLogin"

	var (
		credentials AgentLoginCreds
		agent       models.Agent
		validKey    bool
		token       string
		err         error
	)

	rules := govalidator.MapData{
		"key":      []string{"required"},
		"password": []string{"required"},
	}

	validate := validateRequestPayload(c, rules, &credentials)
	if validate != nil {
		NLog("warning", LogTag, fmt.Sprintf("error validation : %v", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusBadRequest, validate, "Gagal login")
	}

	// check if theres record
	validKey = asira.App.DB.Where("username = ? AND status = ?", credentials.Key, "active").Find(&agent).RecordNotFound()

	if !validKey { // check the password

		err = bcrypt.CompareHashAndPassword([]byte(agent.Password), []byte(credentials.Password))
		if err != nil {
			NLog("error", LogTag, fmt.Sprintf("error password : %v username : %v", err, credentials.Key), c.Get("user").(*jwt.Token), "", true, "")

			return returnInvalidResponse(http.StatusOK, err, "Gagal Login")
		}

		//set role
		tokenrole := "agent"
		token, err = createJwtToken(strconv.FormatUint(agent.ID, 10), tokenrole)
		if err != nil {
			NLog("error", LogTag, fmt.Sprintf("error generating token  : %v ", err), c.Get("user").(*jwt.Token), "", true, "")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat token")
		}
	} else {
		NLog("error", LogTag, fmt.Sprintf("error login  : %v username : %v", err, credentials.Key), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusOK, "", "Gagal Login")
	}

	NLog("event", LogTag, fmt.Sprintf("success login  : %v", credentials.Key), c.Get("user").(*jwt.Token), "", true, "")

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}

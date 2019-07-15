package handlers

import (
	"asira/asira"
	"asira/models"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	BorrowerLoginCreds struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
)

// borrower login, borrower can choose either login with email / phone
func BorrowerLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		credentials BorrowerLoginCreds
		loginType   string
		borrower    models.Borrower
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
		return returnInvalidResponse(http.StatusBadRequest, validate, "invalid login")
	}

	emailchecker := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if emailchecker.MatchString(credentials.Key) {
		loginType = "email"
	}

	// check if theres record
	switch loginType {
	default: // default login is using phone number
		validKey = asira.App.DB.Where("phone = ?", credentials.Key).Find(&borrower).RecordNotFound()
		break
	case "email":
		validKey = asira.App.DB.Where("email = ?", credentials.Key).Find(&borrower).RecordNotFound()
		break
	}

	if !validKey { // check the password
		err = bcrypt.CompareHashAndPassword([]byte(borrower.Password), []byte(credentials.Password))
		if err != nil {
			return returnInvalidResponse(http.StatusOK, err, "invalid login")
		}

		tokenrole := "unverified_borrower"
		if borrower.OTPverified {
			tokenrole = "borrower"
		}
		token, err = createJwtToken(strconv.FormatUint(borrower.ID, 10), tokenrole)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}
	} else {
		return returnInvalidResponse(http.StatusOK, "", "invalid login")
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}

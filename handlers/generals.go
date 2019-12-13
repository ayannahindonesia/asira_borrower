package handlers

import (
	"asira_borrower/asira"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	JWTclaims struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		jwt.StandardClaims
	}
)

// general function to validate all kind of api request payload / body
func validateRequestPayload(c echo.Context, rules govalidator.MapData, data interface{}) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	mappedError := v.ValidateJSON()
	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

// general function to validate all kind of api request url query
func validateRequestQuery(c echo.Context, rules govalidator.MapData) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Rules:   rules,
	}

	v := govalidator.New(opts)

	mappedError := v.Validate()

	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

func returnInvalidResponse(httpcode int, details interface{}, message string) error {
	responseBody := map[string]interface{}{
		"message": message,
		"details": details,
	}

	return echo.NewHTTPError(httpcode, responseBody)
}

// self explanation
func createJwtToken(id string, role string) (string, error) {
	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))

	claim := JWTclaims{
		id,
		role,
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Duration(jwtConf["duration"].(int)) * time.Minute).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := rawToken.SignedString([]byte(jwtConf["jwt_secret"].(string)))
	if err != nil {
		return "", err
	}

	return token, nil
}

func isInArrayInt64(id int64, banks []int64) bool {
	exist := false
	for _, val := range banks {
		if val == id {
			exist = true
			break
		}
	}
	return exist
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

package handlers

import (
	"asira_borrower/asira"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	//JWTclaims hold custom data
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

//isBorrowerRegisteredByAgent check is borrower already registered with agent register borrower
func isBorrowerAlreadyRegistered(idcardNumber string) error {
	db := asira.App.DB
	var count int

	//get users based on borrower id
	db = db.Table("borrowers b").
		Select("u.*").
		Joins("INNER JOIN users u ON b.id = u.borrower").
		Where("b.idcard_number = ?", idcardNumber)

	err = db.Count(&count).Error
	fmt.Println("check err & count ", err, count)
	if err != nil || count > 0 {
		return errors.New("Borrower already registered as personal")
	}

	return nil
}

func checkUniqueFields(idcardNumber string, uniques map[string]string) (string, error) {
	var count int
	fieldsFound := ""

	//...check unique
	for key, val := range uniques {
		//init query
		db := asira.App.DB
		db = db.Table("borrowers").Select("*")

		//get users other than idcardNumber...
		db = db.Not("idcard_number", idcardNumber)

		//if field not empty
		if len(val) > 0 || val != "" {
			db = db.Where(fmt.Sprintf("LOWER(%s) = ?", key), strings.ToLower(val))
		} else {
			//skip checking
			continue
		}
		//query count
		err = db.Count(&count).Error
		fmt.Println("check err & count ", err, count)
		if err != nil || count > 0 {
			fieldsFound += key + ", "
		}
	}
	if fieldsFound != "" {
		return fieldsFound, errors.New("data unique already exist")
	}
	return fieldsFound, nil
}

//imageUpload upload to S3 protocol
func imageUploadFormatted(base64Image string) (string, error) {

	unbased, _ := base64.StdEncoding.DecodeString(base64Image)
	filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
	url, err := asira.App.S3.UploadJPEG(unbased, filename)
	if err != nil {
		return "", err
	}

	return url, nil
}

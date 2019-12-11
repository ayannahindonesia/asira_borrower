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

func checkUniqueFields(idcardNumber string, uniques map[string]string) error {
	db := asira.App.DB
	var (
		count    int
		iterator int = 0
		notEmpty int = 0
	)

	//init query
	db = db.Table("borrowers").Select("*")

	//get users other than idcardNumber...
	db = db.Not("idcard_number", idcardNumber)

	//...check unique
	for key, val := range uniques {
		//DONE: security string from gorm
		if iterator > 0 && (len(val) > 0 || val != "") {
			if notEmpty > 0 {
				db = db.Or(fmt.Sprintf("LOWER(%s) = ?", key), strings.ToLower(val))
			} else {
				db = db.Where(fmt.Sprintf("LOWER(%s) = ?", key), strings.ToLower(val))
			}
			notEmpty++
		}
		iterator++
	}
	//query count
	err = db.Count(&count).Error
	fmt.Println("check err & count ", err, count)
	if err != nil || count > 0 {
		return errors.New("Borrower already registered as personal")
	}
	return nil
}

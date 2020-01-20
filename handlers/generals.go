package handlers

import (
	"asira_borrower/asira"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
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
	db = db.Table("borrowers").
		Select("u.*").
		Joins("INNER JOIN users u ON borrowers.id = u.borrower").
		Where("borrowers.idcard_number = ?", idcardNumber).
		Where(generateDeleteCheck("borrowers"))

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

		//additional check for soft delete
		db = db.Where(generateDeleteCheck("borrowers"))

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

		//additional check for soft delete
		db = db.Where(generateDeleteCheck(tableName))

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

//uploadImageS3 upload to S3 protocol
func uploadImageS3Formatted(prefix string, base64Image string) (string, error) {

	unbased, _ := base64.StdEncoding.DecodeString(base64Image)
	filename := prefix + strconv.FormatInt(time.Now().Unix(), 10)
	url, err := asira.App.S3.UploadJPEG(unbased, filename)
	if err != nil {
		return "", err
	}

	return url, nil
}

//deleteImageS3 delete old image
func deleteImageS3(imageURL string) error {
	i := strings.Split(imageURL, "/")
	delImage := i[len(i)-1]
	err = asira.App.S3.DeleteObject(delImage)
	if err != nil {
		log.Printf("failed to delete image %v from s3 bucket", delImage)
		return err
	}
	return nil
}

//encrypt data with AES 256 CFB
func encrypt(text string, passphrase string) (string, error) {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), err
}

//decrypt data with AES 256 CFB
func decrypt(encryptedText string, passphrase string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(encryptedText)

	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("cannot decrypt")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}

//generateDelete return delete condition : "tablename.deleted_at IS NULL"
func generateDeleteCheck(tableName string) string {
	defaultFormat := "%s.deleted_at IS NULL"
	return fmt.Sprintf(defaultFormat, tableName)
}

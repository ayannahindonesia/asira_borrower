package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ayannahindonesia/northstar/lib/northstarlib"

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

const (
	//NLOGMSG for message body
	NLOGMSG = "message"
	//NLOGERR for error info
	NLOGERR = "error"
	//NLOGQUERY for detailed query tracing
	NLOGQUERY = "query"
)

var EnglishToIndonesiaFields map[string]string = map[string]string{
	"Phone":        "Nomor Telpon",
	"IdCardNumber": "KTP",
	"Email":        "Email",
	"TaxIDnumber":  "NPWP",
}

var EnglishToIndonesiaFieldsUnderscored map[string]string = map[string]string{
	"phone":         "Nomor Telpon",
	"idcard_number": "KTP",
	"email":         "Email",
	"taxid_number":  "NPWP",
}

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
func isBorrowerAlreadyRegistered(email string, phone string) error {
	db := asira.App.DB
	var count int

	//get users based on borrower id
	db = db.Table("borrowers").
		Select("u.*").
		Joins("INNER JOIN users u ON borrowers.id = u.borrower").
		Where("borrowers.phone = ?", phone).
		Or("borrowers.email = ?", email).
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
		if idcardNumber != "" || len(idcardNumber) > 0 {
			db = db.Not("idcard_number", idcardNumber)
		}
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
			word, ok := EnglishToIndonesiaFieldsUnderscored[key]
			if !ok {
				fmt.Println(err)
				word = key
			}
			fieldsFound += word + ", "
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
			word, ok := EnglishToIndonesiaFieldsUnderscored[key]
			if !ok {
				fmt.Println(err)
				word = key
			}
			fieldsFound += word + ", "
		}
	}
	if fieldsFound != "" {
		return fieldsFound, errors.New("data unique already exist")
	}
	return fieldsFound, nil
}

//checkPatchFieldsBorrowers update check field other than id and idcard_number
func checkPatchFieldsBorrowers(id uint64, idcard_number string, uniques map[string]string) (string, error) {

	var count int
	fieldsFound := ""
	tableName := "borrowers"
	fieldID := "id"

	//...check unique
	for key, val := range uniques {
		//init query
		db := asira.App.DB
		db = db.Table(tableName).Select(fieldID)

		//get users other than idcardNumber...
		db = db.Not(fieldID, id)

		if idcard_number != "" || len(idcard_number) > 0 {
			db = db.Not("idcard_number", idcard_number)
		}

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
			word, ok := EnglishToIndonesiaFieldsUnderscored[key]
			if !ok {
				fmt.Println(err)
				word = key
			}
			fieldsFound += word + ", "
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
		log.Println(err)
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

// NLog send log to northstar service
func NLog(level string, tag string, message interface{}, jwttoken *jwt.Token, note string, nouser bool, typeUser string) {
	var (
		uid      string
		username string
		err      error
	)

	if !nouser {
		jti, _ := strconv.ParseUint(jwttoken.Claims.(jwt.MapClaims)["jti"].(string), 10, 64)
		if typeUser == "borrower" {
			user := models.Borrower{}
			err = user.FindbyID(jti)
			if err == nil {
				uid = fmt.Sprint(user.ID)
				username = user.Phone
			}
		} else {
			// agent
			user := models.Agent{}
			err = user.FindbyID(jti)
			if err == nil {
				uid = fmt.Sprint(user.ID)
				username = user.Username
			}
		}
	}

	Message, _ := json.Marshal(message)

	if flag.Lookup("test.v") == nil {
		err = asira.App.Northstar.SubmitKafkaLog(northstarlib.Log{
			Level:    level,
			Tag:      tag,
			Messages: string(Message),
			UID:      uid,
			Username: username,
			Note:     note,
		}, "log")
	}

	if err != nil {
		log.Printf("error northstar log : %v", err)
	}
}

// NAudittrail send audit trail log to northstar service
func NAudittrail(ori interface{}, new interface{}, jwttoken *jwt.Token, entity string, entityID string, action string, typeUser string) {
	var (
		uid      string
		username string
		err      error
	)

	jti, _ := strconv.ParseUint(jwttoken.Claims.(jwt.MapClaims)["jti"].(string), 10, 64)
	if typeUser == "borrower" {
		user := models.Borrower{}
		err = user.FindbyID(jti)
		if err == nil {
			uid = fmt.Sprint(user.ID)
			username = user.Phone
		} else {
			uid = "0"
			username = "not found"
		}
	} else {
		// agent
		user := models.Agent{}
		err = user.FindbyID(jti)
		if err == nil {
			uid = fmt.Sprint(user.ID)
			username = user.Username
		} else {
			uid = "0"
			username = "not found"
		}
	}

	oriMarshal, _ := json.Marshal(ori)
	newMarshal, _ := json.Marshal(new)

	if flag.Lookup("test.v") == nil {
		err = asira.App.Northstar.SubmitKafkaLog(northstarlib.Audittrail{
			Client:   asira.App.Northstar.Secret,
			UserID:   uid,
			Username: username,
			Roles:    typeUser,
			Entity:   entity,
			EntityID: entityID,
			Action:   action,
			Original: fmt.Sprintf(`%s`, string(oriMarshal)),
			New:      fmt.Sprintf(`%s`, string(newMarshal)),
		}, "audittrail")
	}

	if err != nil {
		log.Printf("error northstar log : %v", err)
	}
}

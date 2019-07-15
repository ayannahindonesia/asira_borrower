package tests

import (
	"asira/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestRegisterBorrower(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	admintoken := obj.Value("token").String().Raw()

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+admintoken)
	})

	payload := map[string]interface{}{
		"fullname":              "test fullname",
		"gender":                "M",
		"idcard_number":         "6142356789123",
		"email":                 "mail@domain.com",
		"birthday":              "1999-12-31T00:00:00.000Z",
		"birthplace":            "jakarta",
		"last_education":        "s1",
		"mother_name":           "amamama",
		"phone":                 "08123456789",
		"marriage_status":       "signle",
		"dependants":            0,
		"address":               "sesame st.",
		"province":              "yia ling",
		"city":                  "jakarta",
		"neighbour_association": "001",
		"hamlets":               "002",
		"subdistrict":           "grogol",
		"urban_village":         "jelambar",
		"home_ownership":        "private owned",
		"lived_for":             24,
		"occupation":            "software engineer",
		"employer_name":         "ayannah",
		"employer_address":      "multivision tower",
		"department":            "it",
		"been_workingfor":       3,
		"employer_number":       "08123456789191",
		"monthly_income":        6000000,
		"field_of_work":         "IT",
		"related_personname":    "wuri",
		"related_relation":      "sister",
		"related_phonenumber":   "081247273727",
		"password":              "pass123",
	}

	// test valid response
	obj = auth.POST("/client/register_borrower").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.Keys().Contains("id")

	// test empty payload
	obj = auth.POST("/client/register_borrower").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test unique by registering same data
	obj = auth.POST("/client/register_borrower").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

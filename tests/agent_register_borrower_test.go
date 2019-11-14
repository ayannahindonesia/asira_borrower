package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentBorrower(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	agentToken := getAgentLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+agentToken)
	})

	payload := map[string]interface{}{
		"fullname":              "test fullname",
		"nickname":              "test nickname",
		"gender":                "M",
		"idcard_image":          "iVBORw0KGgoAAAANSUhEUgAAACsAAAAsCAYAAAD8WEF4AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAABFSURBVFhH7c5BDQAwEASh+je9VcHjkkEBb4eUVcoqZZWySlmlrFJWKauUVcoqZZWySlmlrFJWKauUVcoqZZWySlnlUHb7I0d0JGoj43wAAAAASUVORK5CYII=",
		"taxid_image":           "iVBORw0KGgoAAAANSUhEUgAAACsAAAAsCAYAAAD8WEF4AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAABFSURBVFhH7c5BDQAwEASh+je9VcHjkkEBb4eUVcoqZZWySlmlrFJWKauUVcoqZZWySlmlrFJWKauUVcoqZZWySlnlUHb7I0d0JGoj43wAAAAASUVORK5CYII=",
		"idcard_number":         "6142356789123",
		"taxid_number":          "89123",
		"nationality":           "WNI",
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
		"bank":                  1,
		"related_phonenumber":   "081247273727",
		"password":              "pass123",
	}

	// test valid response
	obj := auth.POST("/agent/register_borrower").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.Keys().Contains("id")

	// test empty payload
	obj = auth.POST("/agent/register_borrower").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test unique by registering same data
	obj = auth.POST("/agent/register_borrower").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

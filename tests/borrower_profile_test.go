package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestBorrowerGetProfile(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	// valid response
	obj := auth.GET("/borrower/profile").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// wrong token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrongtoken")
	})
	auth.GET("/borrower/profile").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

func TestBorrowerPatchProfile(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	data := map[string]interface{}{
		"idcard_number":         "61289137218",
		"taxid_number":          "1928374650",
		"mother_name":           "A",
		"address":               "Hiu Putih A",
		"city":                  "Palangkaraya",
		"province":              "yia ling",
		"neighbour_association": "001",
		"hamlets":               "002",
		"subdistrict":           "grogol",
		"urban_village":         "jelambar",
		"occupation":            "dev",
		"bank":                  map[string]interface{}{"Int64": 1, "Valid": true},
	}
	obj := auth.PATCH("/borrower/profile").WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Value("province").Equal("yia ling")

	//test not valid bank id
	data["bank"] = map[string]interface{}{"Int64": 1000, "Valid": true}
	auth.PATCH("/borrower/profile").WithJSON(data).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
}

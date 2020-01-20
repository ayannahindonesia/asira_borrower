package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentBorrowerGetProfile(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	agentToken := getAgentLoginToken(e, auth, "2")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+agentToken)
	})

	// valid response
	obj := auth.GET("/agent/borrower/3").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("monthly_income").ValueEqual("monthly_income", 5000000)

	// wrong token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrongtoken")
	})
	auth.GET("/agent/profile").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

func TestAgentBorrowerProfile(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	agentToken := getAgentLoginToken(e, auth, "2")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+agentToken)
	})

	// valid response
	data := map[string]interface{}{
		"phone":          "081234567899",
		"email":          "agentZ@mail.com",
		"idcard_number":  "614235678912",
		"taxid_number":   "89123128237",
		"monthly_income": 6500000,
	}
	obj := auth.PATCH("/agent/borrower/3").WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("monthly_income").ValueEqual("monthly_income", 6500000)
	//must not updated
	obj.ContainsKey("other_income").ValueEqual("other_income", 2000000)

	//already exist
	data = map[string]interface{}{
		"idcard_number": "9666123451234566689",
	}
	auth.PATCH("/agent/borrower/3").
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// wrong token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrongtoken")
	})
	auth.GET("/agent/profile").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

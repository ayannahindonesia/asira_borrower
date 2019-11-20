package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentChecksBorrower(t *testing.T) {
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

		"idcard_number": "9876123451234566689",
		"phone":         "+629812",
		"email":         "ad@gmail.com",
		"taxid_number":  "0987654321234567890",
	}
	// test valid response
	obj := auth.POST("/agent/checks_borrower").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("status").ValueEqual("status", true)

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer invalid token")
	})
	// test invalid bearer token
	obj = auth.POST("/agent/checks_borrower").
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

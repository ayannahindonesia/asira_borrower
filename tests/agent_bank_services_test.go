package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentBankServicesGet(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getAgentLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	obj := auth.GET("/agent/bank_services/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 2)

	// // valid response of loan details
	// obj = auth.GET("/agent/bank_services/1").
	// 	Expect().
	// 	Status(http.StatusOK).JSON().Object()
	// obj.ContainsKey("id").ValueEqual("id", 1)
	// // loan id not found
	// obj = auth.GET("/agent/bank_services/99").
	// 	Expect().
	// 	Status(http.StatusForbidden).JSON().Object()
}

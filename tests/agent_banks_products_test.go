package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentBankProductsGet(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getAgentLoginToken(e, auth, "2")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	///all products owned by agent's banks
	obj := auth.GET("/agent/bank_products").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 2)

	// valid response of bank_services for bankID = 1
	obj = auth.GET("/agent/bank_products").
		WithQuery("service_id", 1).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 2)

	// valid response of bank_services for bank_id = 1 dan service_id = 1
	obj = auth.GET("/agent/bank_products").
		WithQuery("service_id", 1).
		WithQuery("product_id", 1).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)

	// banks not in agent banks
	obj = auth.GET("/agent/bank_products").
		WithQuery("service_id", 99).
		Expect().
		Status(http.StatusNotFound).JSON().Object()

	// test not found
	obj = auth.GET("/agent/bank_products").
		WithQuery("service_id", 1).
		WithQuery("product_id", 999).
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

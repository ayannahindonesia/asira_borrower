package tests

import (
	"asira_borrower/router"
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
		"fullname": "test fullname",
		"email":    "mail@domain.com",
		"phone":    "08123456789",
		"password": "pass123",
		"otp_code": "888999",
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
		Status(http.StatusInternalServerError).JSON().Object()
}

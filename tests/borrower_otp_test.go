package tests

import (
	"asira/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestBorrowerOTP(t *testing.T) {
	RebuildData()

	token := "Basic YW5kcm9rZXk6YW5kcm9zZWNyZXQ="
	api := router.NewBorrower()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+token)
	})

	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	admintoken := obj.Value("token").String().Raw()

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+admintoken)
	})

	payload := map[string]interface{}{
		"key":      "081234567891",
		"password": "password",
	}

	// test valid response
	obj = auth.POST("/client/borrower_login").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	borrowertoken := obj.Value("token").String().Raw()

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	// test dont have access yet
	obj = auth.GET("/borrower/profile").
		Expect().
		Status(http.StatusForbidden).JSON().Object()

	// valid response
	payload = map[string]interface{}{
		"phone": "081234567891",
	}
	obj = auth.POST("/unverified_borrower/otp_request").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()

	// invalid verify
	payload = map[string]interface{}{
		"phone":    "081234567891",
		"otp_code": "123456",
	}
	obj = auth.POST("/unverified_borrower/otp_verify").WithJSON(payload).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()
	// valid verify
	payload = map[string]interface{}{
		"phone":    "081234567891",
		"otp_code": "888999",
	}
	obj = auth.POST("/unverified_borrower/otp_verify").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
}

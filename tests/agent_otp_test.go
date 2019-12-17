package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentOTP(t *testing.T) {
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

	// test valid response
	payload := map[string]interface{}{
		"phone": "081234567890",
	}
	auth.POST("/agent/otp_request/3").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()

	//valid  verify
	payload = map[string]interface{}{
		"phone":    "081234567890",
		"otp_code": "888999",
	}
	obj := auth.POST("/agent/otp_verify/3").
		WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("message").ValueEqual("message", "OTP Verified")

	// invalid verify
	payload = map[string]interface{}{
		"phone":    "081234567890",
		"otp_code": "123456",
	}
	auth.POST("/agent/otp_verify/3").WithJSON(payload).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	// test invalid bearer token
	auth.POST("/agent/otp_verify/999").
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

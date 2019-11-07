package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestFCMTokenUpdate(t *testing.T) {
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

	payload := map[string]interface{}{
		"fcm_token": "aabbcc",
	}

	// test valid response
	auth.PATCH("/borrower/fcm_token_update").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test invalid client token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrongtoken")
	})
	auth.PATCH("/borrower/fcm_token_update").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

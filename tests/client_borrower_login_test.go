package tests

import (
	"asira/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestBorrowerLogin(t *testing.T) {
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
		"key":      "081234567890",
		"password": "password",
	}

	// test valid response
	obj = auth.POST("/client/borrower_login").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("expires_in", "token")

	// test invalid empty body
	obj = auth.POST("/client/borrower_login").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	// test invalid client token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrongtoken")
	})
	obj = auth.POST("/client/borrower_login").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

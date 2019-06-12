package tests

import (
	"kayacredit/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestClientLogin(t *testing.T) {
	token := "YW5kcm9rZXk6YW5kcm9zZWNyZXQ="
	api := router.NewBorrower()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", token)
	})

	// test valid token
	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("expires_in", "token")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "thisisinvalidtoken")
	})

	// test invalid token
	obj = auth.GET("/clientauth").
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()

	// test without token
	e.GET("/clientauth").
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

package tests

import (
	"asira/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestClientLogin(t *testing.T) {
	api := router.NewBorrower()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
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

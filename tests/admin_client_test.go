package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestClientConfig(t *testing.T) {
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
		"name": "android_client",
		"key":  "android",
		"role": "client",
	}

	// test valid response
	obj = auth.POST("/admin/client_config").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("name", "secret")

	// secret :=
	// // test invalid empty body
	// obj = auth.POST("/client/borrower_login").WithJSON(map[string]interface{}{}).
	// 	Expect().
	// 	Status(http.StatusBadRequest).JSON().Object()

	// // test invalid client token
	// auth = e.Builder(func(req *httpexpect.Request) {
	// 	req.WithHeader("Authorization", "Bearer wrongtoken")
	// })
	// obj = auth.POST("/client/borrower_login").WithJSON(map[string]interface{}{}).
	// 	Expect().
	// 	Status(http.StatusUnauthorized).JSON().Object()
}

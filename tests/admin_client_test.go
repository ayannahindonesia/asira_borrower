package tests

import (
	"asira_borrower/router"
	"encoding/base64"
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

	secret := obj.Value("secret").String().Raw()
	token := "android" + ":" + secret
	var encodedString = base64.StdEncoding.EncodeToString([]byte(token))

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+encodedString)
	})

	obj = auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

}

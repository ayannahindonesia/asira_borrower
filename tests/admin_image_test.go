package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestGetAdminImageString(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// changed to these ==============================
	// auth := e.Builder(func(req *httpexpect.Request) {
	// 	req.WithHeader("Authorization", "Basic "+adminBasicToken)
	// })

	// adminToken := getLenderAdminToken(e, auth)
	// ===============================================

	// These should be changed to... ==================
	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	admintoken := obj.Value("token").String().Raw()
	// ================================================

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+admintoken)
	})

	obj = auth.GET("/admin/image/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)

	// not found
	auth.GET("/admin/image/9999").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

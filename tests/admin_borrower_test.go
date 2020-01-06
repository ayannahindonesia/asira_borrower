package tests

import (
	"testing"
)

func TestBorrowerGetAll(t *testing.T) {
	// 	RebuildData()

	// 	api := router.NewRouter()

	// 	server := httptest.NewServer(api)

	// 	defer server.Close()

	// 	e := httpexpect.New(t, server.URL)

	// 	auth := e.Builder(func(req *httpexpect.Request) {
	// 		req.WithHeader("Authorization", "Basic "+adminBasicToken)
	// 	})

	// 	obj := auth.GET("/clientauth").
	// 		Expect().
	// 		Status(http.StatusOK).JSON().Object()

	// 	admintoken := obj.Value("token").String().Raw()

	// 	auth = e.Builder(func(req *httpexpect.Request) {
	// 		req.WithHeader("Authorization", "Bearer "+admintoken)
	// 	})
	// 	// valid response of borrowers
	// 	obj = auth.GET("/admin/borrower").
	// 		Expect().
	// 		Status(http.StatusOK).JSON().Object()
	// }

	// func TestBorrowerGetDetails(t *testing.T) {
	// 	RebuildData()

	// 	api := router.NewRouter()

	// 	server := httptest.NewServer(api)

	// 	defer server.Close()

	// 	e := httpexpect.New(t, server.URL)

	// 	auth := e.Builder(func(req *httpexpect.Request) {
	// 		req.WithHeader("Authorization", "Basic "+adminBasicToken)
	// 	})
	// 	obj := auth.GET("/clientauth").
	// 		Expect().
	// 		Status(http.StatusOK).JSON().Object()

	// 	admintoken := obj.Value("token").String().Raw()

	// 	auth = e.Builder(func(req *httpexpect.Request) {
	// 		req.WithHeader("Authorization", "Bearer "+admintoken)
	// 	})
	// 	// valid response of loan details
	// 	obj = auth.GET("/admin/borrower/1").
	// 		Expect().
	// 		Status(http.StatusOK).JSON().Object()
	// 	obj.ContainsKey("id").ValueEqual("id", 1)
	// 	// loan id not found
	// 	obj = auth.GET("/admin/borrower/99").
	// 		Expect().
	// 		Status(http.StatusNotFound).JSON().Object()
}

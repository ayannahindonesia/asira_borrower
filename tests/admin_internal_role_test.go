package tests

import (
	"testing"
)

func TestInternalRoleList(t *testing.T) {
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

	// 	// valid response
	// 	obj = auth.GET("/admin/internal_role").
	// 		Expect().
	// 		Status(http.StatusOK).JSON().Object()
	// 	obj.ContainsKey("total_data").NotEmpty()

	// }

	// func TestNewInternalRole(t *testing.T) {
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

	// 	payload := map[string]interface{}{
	// 		"name":        "Admin 2",
	// 		"system":      "Core",
	// 		"description": "Admin Cuy",
	// 		"status":      true,
	// 	}

	// 	// normal scenario
	// 	obj = auth.POST("/admin/internal_role").WithJSON(payload).
	// 		Expect().
	// 		Status(http.StatusCreated).JSON().Object()
	// 	obj.ContainsKey("name").ValueEqual("name", "Admin 2")

	// 	// test invalid
	// 	payload = map[string]interface{}{
	// 		"name": "",
	// 	}
	// 	auth.POST("/admin/internal_role").WithJSON(payload).
	// 		Expect().
	// 		Status(http.StatusUnprocessableEntity).JSON().Object()
	// }

	// func TestGetInternalRoleByID(t *testing.T) {
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

	// 	// valid response
	// 	obj = auth.GET("/admin/internal_role/1").
	// 		Expect().
	// 		Status(http.StatusOK).JSON().Object()
	// 	obj.ContainsKey("id").ValueEqual("id", 1)

	// 	// not found
	// 	auth.GET("/admin/internal_role/9999").
	// 		Expect().
	// 		Status(http.StatusNotFound).JSON().Object()
	// }

	// func TestInternalRolePatch(t *testing.T) {
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

	// 	payload := map[string]interface{}{
	// 		"name": "Test Admin",
	// 	}

	// 	// valid response
	// 	obj = auth.PATCH("/admin/internal_role/1").WithJSON(payload).
	// 		Expect().
	// 		Status(http.StatusOK).JSON().Object()
	// 	obj.ContainsKey("name").ValueEqual("name", "Test Admin")

	// 	// test invalid token
	// 	auth = e.Builder(func(req *httpexpect.Request) {
	// 		req.WithHeader("Authorization", "Bearer wrong token")
	// 	})
	// 	auth.PATCH("/admin/internal_role/1").WithJSON(payload).
	// 		Expect().
	// 		Status(http.StatusUnauthorized).JSON().Object()
}

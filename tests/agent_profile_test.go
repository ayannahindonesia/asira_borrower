package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentGetProfile(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	agentToken := getAgentLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+agentToken)
	})

	// valid response
	obj := auth.GET("/agent/profile").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("username").ValueEqual("username", "agentJ")

	// wrong token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrongtoken")
	})
	auth.GET("/agent/profile").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

func TestAgentPatchProfile(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	agentToken := getAgentLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+agentToken)
	})

	// valid response
	data := map[string]interface{}{
		"phone": "081234567899",
		"email": "agentZ@mail.com",
		"image": "iVBORw0KGgoAAAANSUhEUgAAACsAAAAsCAYAAAD8WEF4AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAABFSURBVFhH7c5BDQAwEASh+je9VcHjkkEBb4eUVcoqZZWySlmlrFJWKauUVcoqZZWySlmlrFJWKauUVcoqZZWySlnlUHb7I0d0JGoj43wAAAAASUVORK5CYII=",
	}
	obj := auth.PATCH("/agent/profile").WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Value("phone").Equal("081234567899")
	obj.Value("email").Equal("agentZ@mail.com")

	obj = auth.GET("/agent/profile").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("username").ValueEqual("username", "agentJ")

	// wrong token
	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer wrongtoken")
	})
	auth.GET("/agent/profile").WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object()
}

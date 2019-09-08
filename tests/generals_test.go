package tests

import (
	"asira_borrower/asira"
	"asira_borrower/migration"
	"fmt"
	"net/http"
	"os"

	"github.com/gavv/httpexpect"
)

var (
	clientBasicToken string = "YW5kcm9rZXk6YW5kcm9zZWNyZXQ="
)

func init() {
	// restrict test to development environment only.
	if asira.App.ENV != "development" {
		fmt.Printf("test aren't allowed in %s environment.", asira.App.ENV)
		os.Exit(1)
	}
}

func RebuildData() {
	migration.Truncate([]string{"all"})
	migration.TestSeed()
}

func getBorrowerLoginToken(e *httpexpect.Expect, auth *httpexpect.Expect, borrower_id string) string {
	obj := auth.GET("/clientauth").
		Expect().
		Status(http.StatusOK).JSON().Object()

	admintoken := obj.Value("token").String().Raw()

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+admintoken)
	})

	var payload map[string]interface{}
	switch borrower_id {
	case "1":
		payload = map[string]interface{}{
			"key":      "081234567890",
			"password": "password",
		}
	case "2":
		payload = map[string]interface{}{
			"key":      "081234567891",
			"password": "password",
		}
	}

	obj = auth.POST("/client/borrower_login").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()

	return obj.Value("token").String().Raw()
}

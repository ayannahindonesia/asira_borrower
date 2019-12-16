package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestAgentLoanGet(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	agenttoken := getAgentLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+agenttoken)
	})

	// valid response of paged loan history
	obj := auth.GET("/agent/loan").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 3)

	obj = auth.GET("/agent/loan").
		WithQuery("status", "approved").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 1)

	// valid response of loan details
	obj = auth.GET("/agent/loan/5/details").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("status").ValueEqual("status", "approved")
	// loan id not found
	obj = auth.GET("/agent/loan/99/details").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestAgentLoanApply(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	agenttoken := getAgentLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+agenttoken)
	})

	payload := map[string]interface{}{
		"borrower":          4,
		"installment":       6,
		"loan_amount":       6000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           1,
	}

	// valid response
	obj := auth.POST("/agent/loan").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("loan_intention").ValueEqual("loan_intention", "Pendidikan")

	obj = auth.GET("/agent/loan/7/details").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("loan_intention").ValueEqual("loan_intention", "Pendidikan")

	// test validation
	payload = map[string]interface{}{
		"installment": "6",
		"loan_amount": "5000000",
	}
	auth.POST("/agent/loan").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
	payload = map[string]interface{}{
		"installment":       "6",
		"loan_amount":       "5000000",
		"loan_intention":    "not valid",
		"intention_details": "the details",
		"product":           1,
	}
	auth.POST("/agent/loan").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	payload = map[string]interface{}{
		"installment":       6,
		"loan_amount":       5000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           99,
	}
	auth.POST("/agent/loan").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test otp
	auth.GET("/agent/loan/7/otp").
		Expect().
		Status(http.StatusOK).JSON().Object()
	// test otp verify
	payload = map[string]interface{}{
		"otp_code": "888999",
	}
	auth.POST("/agent/loan/7/verify").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	// second time should be invalid because loan is already verified
	auth.POST("/agent/loan/7/verify").WithJSON(payload).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()
}

package tests

import (
	"asira_borrower/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestBorrowerLoanGet(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	// valid response of paged loan history
	obj := auth.GET("/borrower/loan").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueEqual("total_data", 3)

	// valid response of loan details
	obj = auth.GET("/borrower/loan/1/details").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)
	// loan id not found
	obj = auth.GET("/borrower/loan/99/details").
		Expect().
		Status(http.StatusNotFound).JSON().Object()
}

func TestBorrowerLoanApply(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	payload := map[string]interface{}{
		"installment":       6,
		"loan_amount":       5000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           1,
	}

	// valid response
	obj := auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("loan_intention").ValueEqual("loan_intention", "Pendidikan")

	// test validation
	payload = map[string]interface{}{
		"installment": "6",
		"loan_amount": "5000000",
	}
	auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
	payload = map[string]interface{}{
		"installment":       "6",
		"loan_amount":       "5000000",
		"loan_intention":    "not valid",
		"intention_details": "the details",
		"product":           1,
	}
	auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	payload = map[string]interface{}{
		"installment":       6,
		"loan_amount":       5000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           99,
	}
	auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()

	// test otp
	auth.GET("/borrower/loan/7/otp").
		Expect().
		Status(http.StatusOK).JSON().Object()

	// test otp verify
	payload = map[string]interface{}{
		"otp_code": "888999",
	}
	auth.POST("/borrower/loan/7/verify").WithJSON(payload).
		Expect().
		Status(http.StatusOK).JSON().Object()
	// second time should be invalid because loan is already verified
	auth.POST("/borrower/loan/7/verify").WithJSON(payload).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()
}

func TestLoanCalculationFormulaFlat(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	// type flat, interest annual 5%
	payload := map[string]interface{}{
		"installment":       6,
		"loan_amount":       5000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           1,
	}

	// valid response
	obj := auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("layaway_plan").ValueEqual("layaway_plan", 875000.0000000001)
	obj.ContainsKey("total_loan").ValueEqual("total_loan", 5250000.000000001)
}

func TestLoanCalculationFormulaFixed(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	// type flat, interest annual 5%
	payload := map[string]interface{}{
		"installment":       12,
		"loan_amount":       5000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           2,
	}

	// valid response
	obj := auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("layaway_plan").ValueEqual("layaway_plan", 434942.0651053465)
	obj.ContainsKey("total_loan").ValueEqual("total_loan", 5219304.781264158)
}

func TestLoanCalculationFormulaOnetimepay(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	// type flat, interest annual 5%
	payload := map[string]interface{}{
		"installment":       8,
		"loan_amount":       5000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           3,
	}

	// valid response
	obj := auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("layaway_plan").ValueEqual("layaway_plan", 687500)
	obj.ContainsKey("total_loan").ValueEqual("total_loan", 5500000)
}

func TestLoanCalculationFormulaEfektifmenurun(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	borrowertoken := getBorrowerLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+borrowertoken)
	})

	// type flat, interest annual 5%
	payload := map[string]interface{}{
		"installment":       24,
		"loan_amount":       10000000,
		"loan_intention":    "Pendidikan",
		"intention_details": "the details",
		"product":           4,
	}

	// valid response
	obj := auth.POST("/borrower/loan").WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("total_loan").ValueEqual("total_loan", 11250000)
}

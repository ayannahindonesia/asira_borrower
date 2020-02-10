package tests

import (
	"testing"
)

func TestBorrowerOTP(t *testing.T) {
	// RebuildData()

	// api := router.NewRouter()

	// server := httptest.NewServer(api)

	// defer server.Close()

	// e := httpexpect.New(t, server.URL)

	// auth := e.Builder(func(req *httpexpect.Request) {
	// 	req.WithHeader("Authorization", "Basic "+clientBasicToken)
	// })

	// borrowertoken := getBorrowerLoginToken(e, auth, "2")

	// auth = e.Builder(func(req *httpexpect.Request) {
	// 	req.WithHeader("Authorization", "Bearer "+borrowertoken)
	// })

	// // test dont have access yet
	// auth.GET("/borrower/profile").
	// 	Expect().
	// 	Status(http.StatusForbidden).JSON().Object()

	// // valid response
	// payload := map[string]interface{}{
	// 	"phone": "081234567891",
	// }
	// auth.POST("/unverified_borrower/otp_request").WithJSON(payload).
	// 	Expect().
	// 	Status(http.StatusOK).JSON().Object()

	// // invalid verify
	// payload = map[string]interface{}{
	// 	"phone":    "081234567891",
	// 	"otp_code": "123456",
	// }
	// auth.POST("/unverified_borrower/otp_verify").WithJSON(payload).
	// 	Expect().
	// 	Status(http.StatusBadRequest).JSON().Object()
	// // valid verify
	// payload = map[string]interface{}{
	// 	"phone":    "081234567891",
	// 	"otp_code": "888999",
	// }
	// auth.POST("/unverified_borrower/otp_verify").WithJSON(payload).
	// 	Expect().
	// 	Status(http.StatusOK).JSON().Object()
}

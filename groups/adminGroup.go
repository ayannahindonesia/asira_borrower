package groups

import (
	"github.com/labstack/echo"
)

func AdminGroup(e *echo.Echo) {
	// g := e.Group("/admin")
	// middlewares.SetClientJWTmiddlewares(g, "admin")

	// // config info
	// g.GET("/info", handlers.AsiraAppInfo)

	// // Loan Purpose
	// g.GET("/loan_purposes", admin_handlers.LoanPurposeList)
	// g.POST("/loan_purposes", admin_handlers.LoanPurposeNew)
	// g.GET("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDetail)
	// g.PATCH("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposePatch)
	// g.DELETE("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDelete)

}

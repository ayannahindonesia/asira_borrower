package validator

import (
	"fmt"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

type AsiraValidator struct{}

func (a *AsiraValidator) customValidatorRules() {
	govalidator.AddCustomRule("asc_desc", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if strings.ToUpper(val) != "ASC" && strings.ToUpper(val) != "DESC" {
			return fmt.Errorf("The %s field must be contain word: asc or desc", field)
		}
		return nil
	})

	govalidator.AddCustomRule("loan_statuses", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if val != "approved" && val != "rejected" && val != "processing" {
			return fmt.Errorf("The %s field must be contain: approved, rejected, or processing", field)
		}
		return nil
	})
}

package validator

import (
	"fmt"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

type AsiraValidator struct{}

func (a *AsiraValidator) CustomValidatorRules() {
	// validator for pagination
	govalidator.AddCustomRule("asc_desc", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if strings.ToUpper(val) != "ASC" && strings.ToUpper(val) != "DESC" {
			return fmt.Errorf("The %s field must be contain word: asc or desc", field)
		}
		return nil
	})

	// validator for loans
	govalidator.AddCustomRule("loan_statuses", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if val != "approved" && val != "rejected" && val != "processing" {
			return fmt.Errorf("The %s field must be contain either: approved, rejected, or processing", field)
		}
		return nil
	})

	// validator for otp entity types
	govalidator.AddCustomRule("otp_entity_types", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if val != "loan" && val != "borrower" {
			return fmt.Errorf("The %s field must be contain either: loan or borrower", field)
		}
		return nil
	})
}

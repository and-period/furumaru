package codes

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const accountIDString = "^[a-zA-Z0-9_-]*$"

var accountIDRegex = regexp.MustCompile(accountIDString)

func RegisterValidations(v *validator.Validate) error {
	validators := map[string]validator.Func{
		"account_id": validateAccountID,
	}
	for tag, f := range validators {
		if err := v.RegisterValidation(tag, f); err != nil {
			return err
		}
	}

	return nil
}

func validateAccountID(fl validator.FieldLevel) bool {
	return accountIDRegex.MatchString(fl.Field().String())
}

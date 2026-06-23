package dto

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	passwordHasUpper   = regexp.MustCompile(`[A-Z]`)
	passwordHasLower   = regexp.MustCompile(`[a-z]`)
	passwordHasNumber  = regexp.MustCompile(`[0-9]`)
	passwordHasSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-={}\[\]|;:"<>,./?]`)
)

// RegisterValidators registers all custom validators with gin's validator engine.
// Add new validators here — main.go only calls this single function.
func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil
	}

	if err := v.RegisterValidation("password", passwordValidator); err != nil {
		return err
	}

	// Future validators go here:
	// if err := v.RegisterValidation("phone", phoneValidator); err != nil {
	//     return err
	// }

	return nil
}

// passwordValidator checks: 8-72 chars, uppercase, lowercase, number, special char.
func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 || len(password) > 72 {
		return false
	}
	return passwordHasUpper.MatchString(password) &&
		passwordHasLower.MatchString(password) &&
		passwordHasNumber.MatchString(password) &&
		passwordHasSpecial.MatchString(password)
}

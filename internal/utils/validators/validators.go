package validators

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"unicode"

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
func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil
	}

	// Use json tag names in validation error messages instead of struct field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := v.RegisterValidation("password", passwordValidator); err != nil {
		return err
	}

	if err := v.RegisterValidation("phone", phoneValidator); err != nil {
		return err
	}

	return nil
}

// FormatValidationError converts validator errors into human-readable messages.
func FormatValidationError(err error) error {
	// Handle empty request body
	if errors.Is(err, io.EOF) {
		return errors.New("request body is required; please provide the required JSON fields")
	}

	// Handle malformed JSON / syntax errors
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) {
		return errors.New("invalid JSON format in request body")
	}
	if errors.As(err, &unmarshalErr) {
		return fmt.Errorf("invalid type for field '%s': expected %s, got %s", unmarshalErr.Field, unmarshalErr.Type.String(), unmarshalErr.Value)
	}

	// Handle validation errors (missing fields, wrong format, etc.)
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	var msgs []string
	for _, e := range validationErrors {
		field := e.Field()
		tag := e.Tag()
		param := e.Param()

		var msg string
		switch tag {
		case "required":
			msg = fmt.Sprintf("%s is required", field)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address", field)
		case "password":
			msg = "password must be 8-72 characters and contain uppercase, lowercase, number, and special character"
		case "phone":
			msg = "phone must be digits and exactly 11 digits"
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", field, param)
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", field, param)
		default:
			msg = fmt.Sprintf("%s is invalid", field)
		}
		msgs = append(msgs, msg)
	}

	return errors.New(strings.Join(msgs, "; "))
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

// phoneValidator checks and make sure the phone number is == 11
func phoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 11 {
		return false
	}

	for _, t := range phone {
		if !unicode.IsDigit(t) {
			return false
		}

	}

	return true
}

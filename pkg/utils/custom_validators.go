package utils

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// RegisterVietnamesePhoneValidator registers Vietnamese phone number validation
// Accepts formats:
//   - 0912345678 (10 digits starting with 0)
//   - +84912345678 (international format)
//   - 84912345678 (without +)
func RegisterVietnamesePhoneValidator(cv *CustomValidator) error {
	// Vietnamese phone number patterns
	viPhoneRegex := regexp.MustCompile(`^(0|\+84|84)(3[2-9]|5[689]|7[06-9]|8[1-9]|9[0-9])\d{7}$`)

	err := cv.validator.RegisterValidation("vi_phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return true // Let 'required' tag handle empty values
		}
		return viPhoneRegex.MatchString(phone)
	})

	if err != nil {
		return err
	}

	// Register translation for Vietnamese phone validator
	return cv.validator.RegisterTranslation(
		"vi_phone",
		cv.translator,
		func(ut ut.Translator) error {
			return ut.Add("vi_phone", "{0} must be a valid Vietnamese phone number", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("vi_phone", fe.Field())
			return t
		},
	)
}

// RegisterAllCustomValidators registers all custom validators at once
func (cv *CustomValidator) RegisterAllCustomValidators() error {
	// Register Vietnamese phone validator
	if err := RegisterVietnamesePhoneValidator(cv); err != nil {
		return err
	}

	return nil
}

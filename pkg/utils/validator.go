package utils

import (
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type CustomValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

// NewValidator creates a new validator instance with dependency injection
// No global state - translator is injected
func NewValidator(trans ut.Translator) *CustomValidator {
	// Initialize validator
	validate := validator.New()

	// Register default translations
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	// Register tag name function to use json tags for field names
	// This makes error field names match JSON keys (lowercase)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{
		validator:  validate,
		translator: trans,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// ExtractValidationErrors extracts field-level validation errors using translator
// Returns map[string]string for easier client-side consumption
// Key = field name (lowercase from json tag)
// Value = translated error message
func (cv *CustomValidator) ExtractValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			// Use translator to get localized error message
			// fieldErr.Field() returns json tag name (already lowercase from RegisterTagNameFunc)
			errorsMap[fieldErr.Field()] = fieldErr.Translate(cv.translator)
		}
		return errorsMap
	}

	// If error is not validator.ValidationErrors, return empty map
	return errorsMap
}

// RegisterCustomValidator registers a custom validation rule on this validator instance
// No global state - uses instance translator via dependency injection
func (cv *CustomValidator) RegisterCustomValidator(
	tag string,
	fn validator.Func,
	message string,
) error {
	err := cv.validator.RegisterValidation(tag, fn)
	if err != nil {
		return err
	}

	// Register translation for custom validator
	return cv.validator.RegisterTranslation(
		tag,
		cv.translator, // ← Use instance translator, not global
		func(ut ut.Translator) error {
			return ut.Add(tag, message, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		},
	)
}

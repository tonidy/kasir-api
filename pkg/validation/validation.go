package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator library for input validation
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	validate := validator.New()

	// Register custom translations or validations if needed
	return &Validator{
		validate: validate,
	}
}

// ValidateStruct validates a struct based on validation tags
func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		// Translate validation errors to a more readable format
		return v.translateValidationErrors(err)
	}
	return nil
}

// ValidateField validates a single field
func (v *Validator) ValidateField(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// translateValidationErrors converts validator errors to a user-friendly format
func (v *Validator) translateValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	var errorMessages []string

	for _, e := range validationErrors {
		fieldName := e.Field()
		tag := e.Tag()

		var msg string
		switch tag {
		case "required":
			msg = fmt.Sprintf("%s is required", fieldName)
		case "min":
			minValue := e.Param()
			msg = fmt.Sprintf("%s must be at least %s", fieldName, minValue)
		case "max":
			maxValue := e.Param()
			msg = fmt.Sprintf("%s must be at most %s", fieldName, maxValue)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "len":
			length := e.Param()
			msg = fmt.Sprintf("%s must be exactly %s characters long", fieldName, length)
		case "gt":
			minValue := e.Param()
			msg = fmt.Sprintf("%s must be greater than %s", fieldName, minValue)
		case "gte":
			minValue := e.Param()
			msg = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, minValue)
		case "lt":
			maxValue := e.Param()
			msg = fmt.Sprintf("%s must be less than %s", fieldName, maxValue)
		case "lte":
			maxValue := e.Param()
			msg = fmt.Sprintf("%s must be less than or equal to %s", fieldName, maxValue)
		case "oneof":
			values := e.Param()
			msg = fmt.Sprintf("%s must be one of [%s]", fieldName, values)
		default:
			msg = fmt.Sprintf("%s failed validation: %s", fieldName, tag)
		}

		errorMessages = append(errorMessages, msg)
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
}

// Custom validation functions can be registered here
func (v *Validator) RegisterValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

// GetFieldErrors extracts field-specific validation errors
func (v *Validator) GetFieldErrors(err error) map[string]string {
	fieldErrors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fieldName := strings.ToLower(e.Field())
			fieldErrors[fieldName] = v.getErrorMessage(e)
		}
	}

	return fieldErrors
}

// getErrorMessage gets a human-readable error message for a validation error
func (v *Validator) getErrorMessage(e validator.FieldError) string {
	fieldName := e.Field()
	tag := e.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "min":
		return fmt.Sprintf("%s must be at least %s", fieldName, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", fieldName, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldName)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fieldName, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fieldName, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fieldName, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fieldName, e.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", fieldName, tag)
	}
}

// ValidateStructWithTags validates a struct and returns field-specific errors
func (v *Validator) ValidateStructWithTags(s interface{}) (map[string]string, error) {
	if err := v.validate.Struct(s); err != nil {
		return v.GetFieldErrors(err), err
	}
	return nil, nil
}

// ValidateStructWithResult validates a struct and returns detailed validation result
func (v *Validator) ValidateStructWithResult(s interface{}) *ValidationResult {
	err := v.validate.Struct(s)

	result := &ValidationResult{
		Valid: true,
	}

	if err != nil {
		result.Valid = false
		result.Errors = v.GetFieldErrors(err)
		result.RawError = err
	}

	return result
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid    bool
	Errors   map[string]string
	RawError error
}

// HasErrors returns true if there are validation errors
func (vr *ValidationResult) HasErrors() bool {
	return !vr.Valid
}

// GetError returns the first validation error as a string
func (vr *ValidationResult) GetError() string {
	if vr.Valid || len(vr.Errors) == 0 {
		return ""
	}

	for _, err := range vr.Errors {
		return err
	}

	return "validation failed"
}

// GetErrors returns all validation errors as a slice of strings
func (vr *ValidationResult) GetErrors() []string {
	var errors []string
	for _, err := range vr.Errors {
		errors = append(errors, err)
	}
	return errors
}

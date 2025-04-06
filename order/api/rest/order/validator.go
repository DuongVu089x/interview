package order

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator is a custom validator for echo framework
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new custom validator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate validates the given struct
func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

// Ensure CustomValidator implements echo.Validator interface
var _ echo.Validator = (*CustomValidator)(nil)

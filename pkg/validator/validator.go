package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps go-playground/validator for Echo.
type CustomValidator struct {
	v *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{v: validator.New()}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.v.Struct(i)
}

package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf(
				"field '%s' failed validation for tag '%s'",
				strings.ToLower(err.Field()),
				err.Tag(),
			))
		}
		return errors.New(strings.Join(validationErrors, ", "))
	}
	return nil
}

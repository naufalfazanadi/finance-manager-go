package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := &Validator{
		validate: validator.New(),
	}

	// Register custom validators
	v.validate.RegisterValidation("datetime", v.validateDateTime)

	return v
}

// validateDateTime validates date time format
func (v *Validator) validateDateTime(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	layout := fl.Param()

	if dateStr == "" {
		return false
	}

	_, err := time.Parse(layout, dateStr)
	return err == nil
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "datetime":
				validationErrors = append(validationErrors, fmt.Sprintf(
					"field '%s' must be a valid date in format %s",
					strings.ToLower(err.Field()),
					err.Param(),
				))
			default:
				validationErrors = append(validationErrors, fmt.Sprintf(
					"field '%s' failed validation for tag '%s'",
					strings.ToLower(err.Field()),
					err.Tag(),
				))
			}
		}
		return errors.New(strings.Join(validationErrors, ", "))
	}
	return nil
}

// ParseJSONStrict parses JSON with DisallowUnknownFields to reject extra fields
func (v *Validator) ParseJSONStrict(c *fiber.Ctx, dst interface{}) error {
	body := c.Body()
	if len(body) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "empty request body")
	}

	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

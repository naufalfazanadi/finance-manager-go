package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// Types and Structs
// ============================================================================

// Validator wraps the go-playground validator with custom validation rules
type Validator struct {
	validate *validator.Validate
}

// FileValidation represents file validation configuration
type FileValidation struct {
	MaxSize      int64    // Maximum file size in bytes
	AllowedTypes []string // Allowed MIME types
	Required     bool     // Whether the file is required
}

// FileValidationResult represents the result of file validation
type FileValidationResult struct {
	Valid bool   // Whether the file is valid
	Error string // Error message if invalid
}

// ============================================================================
// Constructor
// ============================================================================

// New creates a new validator instance with custom validation rules
func New() *Validator {
	v := &Validator{
		validate: validator.New(),
	}

	// Register custom validators
	v.validate.RegisterValidation("datetime", v.validateDateTime)
	v.validate.RegisterValidation("strongpassword", v.validateStrongPassword)

	return v
}

// ============================================================================
// Core Validation Methods
// ============================================================================

// Validate validates a struct using the registered validation rules
// Returns formatted, human-readable error messages
func (v *Validator) Validate(i interface{}) error {
	err := v.validate.Struct(i)
	if err != nil {
		formattedErrors := v.formatValidationErrors(err)
		return errors.New(strings.Join(formattedErrors, "; "))
	}
	return nil
}

// ValidateFile validates a file upload based on the provided configuration
func (v *Validator) ValidateFile(file *multipart.FileHeader, validation FileValidation) FileValidationResult {
	// Check if file is required
	if file == nil {
		if validation.Required {
			return FileValidationResult{
				Valid: false,
				Error: "File is required",
			}
		}
		return FileValidationResult{Valid: true}
	}

	// Check file size
	if file.Size > validation.MaxSize {
		return FileValidationResult{
			Valid: false,
			Error: fmt.Sprintf("File size must not exceed %d bytes", validation.MaxSize),
		}
	}

	// Check file type if types are specified
	if len(validation.AllowedTypes) > 0 {
		contentType := v.detectFileContentType(file)
		if contentType == "unknown" {
			return FileValidationResult{
				Valid: false,
				Error: "Unable to determine file type",
			}
		}

		if !v.isContentTypeAllowed(contentType, validation.AllowedTypes) {
			return FileValidationResult{
				Valid: false,
				Error: fmt.Sprintf("Invalid file type. Expected %s, got %s",
					strings.Join(validation.AllowedTypes, " or "), contentType),
			}
		}
	}

	return FileValidationResult{Valid: true}
}

// ============================================================================
// JSON Parsing Methods
// ============================================================================

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

// ParseAndValidate parses with strict JSON validation and struct validation
func (v *Validator) ParseAndValidate(c *fiber.Ctx, dst interface{}) error {
	if err := v.ParseJSONStrict(c, dst); err != nil {
		return err
	}
	return v.Validate(dst)
}

// ============================================================================
// Form Parsing Methods
// ============================================================================

// ParseFormStrict parses multipart form data with strict field validation
func (v *Validator) ParseFormStrict(c *fiber.Ctx, dst interface{}) error {
	// Parse the multipart form
	if err := c.BodyParser(dst); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse form data: "+err.Error())
	}

	// Handle file uploads manually
	if err := v.parseFileFields(c, dst); err != nil {
		return err
	}

	// Get all form values to check for unknown fields
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse multipart form: "+err.Error())
	}

	// Get allowed field names from the struct tags
	allowedFields := v.getAllowedFormFields(dst)

	// Check for unknown fields in form values
	for fieldName := range form.Value {
		if !v.contains(allowedFields, fieldName) {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Unknown field: %s", fieldName))
		}
	}

	// Check for unknown fields in form files
	for fieldName := range form.File {
		if !v.contains(allowedFields, fieldName) {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Unknown file field: %s", fieldName))
		}
	}

	return nil
}

// ParseFormAndValidate parses multipart form data with strict field validation and struct validation
func (v *Validator) ParseFormAndValidate(c *fiber.Ctx, dst interface{}) error {
	if err := v.ParseFormStrict(c, dst); err != nil {
		return err
	}
	return v.Validate(dst)
}

// ============================================================================
// Custom Validation Functions
// ============================================================================

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

// validateStrongPassword validates password strength:
// - At least 1 uppercase letter
// - At least 1 number
// - At least 1 special character/symbol
func (v *Validator) validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if password == "" {
		return false
	}

	patterns := []string{
		`[A-Z]`, // uppercase letter
		`[0-9]`, // number
		`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`, // special character
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, password)
		if !matched {
			return false
		}
	}

	return true
}

// ============================================================================
// Helper Methods
// ============================================================================

// formatValidationErrors formats validation errors into readable messages
func (v *Validator) formatValidationErrors(err error) []string {
	var validationErrors []string

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	for _, err := range validationErrs {
		fieldName := v.getJSONFieldName(err.StructNamespace(), err.Field())
		message := v.getValidationErrorMessage(err, fieldName)
		validationErrors = append(validationErrors, message)
	}

	return validationErrors
}

// getValidationErrorMessage returns a user-friendly error message for a validation error
func (v *Validator) getValidationErrorMessage(err validator.FieldError, fieldName string) string {
	switch err.Tag() {
	case "datetime":
		return fmt.Sprintf("field '%s' must be a valid date in format %s", fieldName, err.Param())
	case "strongpassword":
		return fmt.Sprintf("field '%s' must contain at least 1 uppercase letter, 1 number, and 1 special character", fieldName)
	case "required":
		return fmt.Sprintf("field '%s' is required", fieldName)
	case "email":
		return fmt.Sprintf("field '%s' must be a valid email address", fieldName)
	case "min":
		return fmt.Sprintf("field '%s' must be at least %s characters long", fieldName, err.Param())
	case "max":
		return fmt.Sprintf("field '%s' must not exceed %s characters", fieldName, err.Param())
	default:
		return fmt.Sprintf("field '%s' failed validation for tag '%s'", fieldName, err.Tag())
	}
}

// getJSONFieldName returns the JSON field name for a struct field
func (v *Validator) getJSONFieldName(structNamespace, fieldName string) string {
	// TODO: This could be enhanced to parse struct tags if needed
	return strings.ToLower(fieldName)
}

// getAllowedFormFields extracts form field names from struct tags
func (v *Validator) getAllowedFormFields(dst interface{}) []string {
	var allowedFields []string

	t := reflect.TypeOf(dst)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check form tag first
		if formTag := field.Tag.Get("form"); formTag != "" && formTag != "-" {
			fieldName := strings.Split(formTag, ",")[0]
			if fieldName != "" {
				allowedFields = append(allowedFields, fieldName)
			}
			continue
		}

		// Fallback to json tag
		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			fieldName := strings.Split(jsonTag, ",")[0]
			if fieldName != "" {
				allowedFields = append(allowedFields, fieldName)
			}
		}
	}

	return allowedFields
}

// detectFileContentType detects the content type from file
func (v *Validator) detectFileContentType(file *multipart.FileHeader) string {
	fileReader, err := file.Open()
	if err != nil {
		return "unknown"
	}
	defer fileReader.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = fileReader.Read(buffer)
	if err != nil {
		return "unknown"
	}

	return v.detectContentType(buffer)
}

// detectContentType detects the content type from file bytes
func (v *Validator) detectContentType(data []byte) string {
	// Check for JPEG
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}

	// Check for PNG
	if len(data) >= 8 &&
		data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "image/png"
	}

	// Check for GIF
	if len(data) >= 6 &&
		data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 &&
		data[3] == 0x38 && (data[4] == 0x37 || data[4] == 0x39) && data[5] == 0x61 {
		return "image/gif"
	}

	return "unknown"
}

// isContentTypeAllowed checks if the content type is in the allowed list
func (v *Validator) isContentTypeAllowed(contentType string, allowedTypes []string) bool {
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

// parseFileFields manually parses file fields from multipart form and assigns them to struct fields
func (v *Validator) parseFileFields(c *fiber.Ctx, dst interface{}) error {
	// Get the struct value and type
	val := reflect.ValueOf(dst)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	// Get the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return nil // No multipart form, skip file parsing
	}

	// Iterate through struct fields to find file fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !fieldVal.CanSet() {
			continue
		}

		// Check if this is a file field (*multipart.FileHeader)
		if fieldVal.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
			// Get the form tag name
			formTag := field.Tag.Get("form")
			if formTag == "" || formTag == "-" {
				continue
			}

			// Get the field name from the tag
			fieldName := strings.Split(formTag, ",")[0]

			// Check if a file was uploaded for this field
			if files, exists := form.File[fieldName]; exists && len(files) > 0 {
				// Set the first file to the field
				fieldVal.Set(reflect.ValueOf(files[0]))
			}
		}
	}

	return nil
}

// contains checks if a string slice contains a specific string
func (v *Validator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

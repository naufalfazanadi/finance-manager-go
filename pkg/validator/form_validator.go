package validator

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/pkg/upload"
)

// FormValidationConfig represents configuration for form validation
type FormValidationConfig struct {
	RequiredFields   []string // List of required form fields
	OptionalFields   []string // List of optional form fields
	FileFields       []string // List of file upload fields
	ValidateFiles    bool     // Whether to validate uploaded files
	AllowEmptyUpdate bool     // For update operations, allow empty values
}

// FormDataResult represents the result of form validation
type FormDataResult struct {
	Data  interface{}                      // The validated DTO
	Files map[string]*multipart.FileHeader // Map of field name to file header
}

// ValidateForm is a general form validation function that can handle any DTO type
//
// Example usage:
//
//	// For a custom DTO
//	req := &dto.CustomRequest{}
//	config := FormValidationConfig{
//	  RequiredFields: []string{"field1", "field2"},
//	  OptionalFields: []string{"field3"},
//	  FileFields:     []string{"document", "image"},
//	  ValidateFiles:  true,
//	}
//	result, err := ValidateForm(c, validator, req, config)
//	if err != nil {
//	  return handleError(err)
//	}
//
//	customReq := result.Data.(*dto.CustomRequest)
//	documentFile := result.Files["document"]
//	imageFile := result.Files["image"]
func (v *Validator) ValidateForm(c *fiber.Ctx, dtoPtr interface{}, config FormValidationConfig) (*FormDataResult, error) {
	// Get the type and value of the DTO
	dtoValue := reflect.ValueOf(dtoPtr)
	if dtoValue.Kind() != reflect.Ptr || dtoValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid DTO type: DTO must be a pointer to struct")
	}

	dtoElem := dtoValue.Elem()
	dtoType := dtoElem.Type()

	// Parse form fields and populate the DTO
	for i := 0; i < dtoType.NumField(); i++ {
		field := dtoType.Field(i)
		fieldValue := dtoElem.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Get the form field name (use json tag or field name)
		formFieldName := getFormFieldName(field)
		formValue := c.FormValue(formFieldName)

		// Set the field value based on its type
		if err := setFieldValue(fieldValue, formValue, field, config.AllowEmptyUpdate); err != nil {
			return nil, fmt.Errorf("invalid %s format: %w", formFieldName, err)
		}
	}

	// Validate the populated DTO
	if err := v.Validate(dtoPtr); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Handle file uploads if configured
	files := make(map[string]*multipart.FileHeader)
	if config.ValidateFiles {
		for _, fileField := range config.FileFields {
			fileHeader, err := c.FormFile(fileField)
			if err == nil && fileHeader != nil {
				// Validate file based on field name
				if err := validateFileByField(fileField, fileHeader); err != nil {
					return nil, fmt.Errorf("file validation failed for %s: %w", fileField, err)
				}
				files[fileField] = fileHeader
			}
		}
	}

	return &FormDataResult{
		Data:  dtoPtr,
		Files: files,
	}, nil
}

// getFormFieldName extracts the form field name from struct field
func getFormFieldName(field reflect.StructField) string {
	// Check for json tag first
	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		// Remove options like omitempty
		if commaIdx := len(jsonTag); commaIdx > 0 {
			for i, r := range jsonTag {
				if r == ',' {
					commaIdx = i
					break
				}
			}
			return jsonTag[:commaIdx]
		}
		return jsonTag
	}

	// Check for form tag
	if formTag := field.Tag.Get("form"); formTag != "" {
		return formTag
	}

	// Convert field name to snake_case
	return toSnakeCase(field.Name)
}

// setFieldValue sets the field value based on the form input and field type
func setFieldValue(fieldValue reflect.Value, formValue string, field reflect.StructField, allowEmpty bool) error {
	// Skip empty values for update operations if allowed
	if allowEmpty && formValue == "" {
		return nil
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(formValue)

	case reflect.Ptr:
		if formValue == "" {
			fieldValue.Set(reflect.Zero(fieldValue.Type()))
			return nil
		}

		// Handle pointer to time.Time for birth_date fields
		if fieldValue.Type() == reflect.TypeOf((*time.Time)(nil)) {
			if parsedTime, err := time.Parse(time.RFC3339, formValue); err == nil {
				fieldValue.Set(reflect.ValueOf(&parsedTime))
			} else {
				return err
			}
		}

	// Add more type handlers as needed
	case reflect.Int, reflect.Int64:
		// Handle integer fields if needed in the future
	}

	return nil
}

// validateFileByField validates files based on the field name
func validateFileByField(fieldName string, fileHeader *multipart.FileHeader) error {
	switch fieldName {
	case "profile_photo":
		return upload.ValidateProfilePhoto(fileHeader)
	// Add more file validation cases as needed
	default:
		// Default file validation or no validation
		return nil
	}
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		if r >= 'A' && r <= 'Z' {
			result = append(result, r-'A'+'a')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

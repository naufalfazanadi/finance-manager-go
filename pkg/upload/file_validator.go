package upload

import (
	"fmt"
	"mime/multipart"
	"strings"
)

const (
	MaxFileSize    = 2 * 1024 * 1024 // 2MB in bytes
	AllowedMimeJPG = "image/jpeg"
	AllowedMimePNG = "image/png"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateProfilePhoto validates the uploaded profile photo
func ValidateProfilePhoto(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return ValidationError{
			Field:   "profile_photo",
			Message: "Profile photo is required",
		}
	}

	// Check file size
	if fileHeader.Size > MaxFileSize {
		return ValidationError{
			Field:   "profile_photo",
			Message: fmt.Sprintf("File size must not exceed %d bytes (2MB)", MaxFileSize),
		}
	}

	// Check file extension
	filename := strings.ToLower(fileHeader.Filename)
	if !strings.HasSuffix(filename, ".jpg") &&
		!strings.HasSuffix(filename, ".jpeg") &&
		!strings.HasSuffix(filename, ".png") {
		return ValidationError{
			Field:   "profile_photo",
			Message: "Only JPG and PNG files are allowed",
		}
	}

	// Open file to check MIME type
	file, err := fileHeader.Open()
	if err != nil {
		return ValidationError{
			Field:   "profile_photo",
			Message: "Failed to open uploaded file",
		}
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return ValidationError{
			Field:   "profile_photo",
			Message: "Failed to read file content",
		}
	}

	// Reset file pointer to beginning
	file.Seek(0, 0)

	// Detect content type
	contentType := detectContentType(buffer)
	if contentType != AllowedMimeJPG && contentType != AllowedMimePNG {
		return ValidationError{
			Field:   "profile_photo",
			Message: fmt.Sprintf("Invalid file type. Expected %s or %s, got %s", AllowedMimeJPG, AllowedMimePNG, contentType),
		}
	}

	return nil
}

// detectContentType detects the content type from file bytes
func DetectContentType(data []byte) string {
	// Check for JPEG
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return AllowedMimeJPG
	}

	// Check for PNG
	if len(data) >= 8 &&
		data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return AllowedMimePNG
	}

	return "unknown"
}

// detectContentType detects the content type from file bytes (internal)
func detectContentType(data []byte) string {
	return DetectContentType(data)
}

// GetFileExtension returns the file extension based on content type
func GetFileExtension(contentType string) string {
	switch contentType {
	case AllowedMimeJPG:
		return ".jpg"
	case AllowedMimePNG:
		return ".png"
	default:
		return ""
	}
}

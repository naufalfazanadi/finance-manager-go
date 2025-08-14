package upload

import "github.com/naufalfazanadi/finance-manager-go/pkg/validator"

// Common file validation configurations that can be reused across the application

var (
	// ProfilePhotoValidation defines validation rules for profile photo uploads
	ProfilePhotoValidation = validator.FileValidation{
		MaxSize:      2 * 1024 * 1024, // 2MB
		AllowedTypes: []string{"image/jpeg", "image/png"},
		Required:     false,
	}

	// RequiredProfilePhotoValidation defines validation rules for required profile photo uploads
	RequiredProfilePhotoValidation = validator.FileValidation{
		MaxSize:      2 * 1024 * 1024, // 2MB
		AllowedTypes: []string{"image/jpeg", "image/png"},
		Required:     true,
	}

	// DocumentValidation defines validation rules for document uploads (PDF, DOC, DOCX)
	DocumentValidation = validator.FileValidation{
		MaxSize:      5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		Required:     false,
	}

	// ImageValidation defines general validation rules for image uploads
	ImageValidation = validator.FileValidation{
		MaxSize:      3 * 1024 * 1024, // 3MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
		Required:     false,
	}

	// AvatarValidation defines validation rules for avatar/small image uploads
	AvatarValidation = validator.FileValidation{
		MaxSize:      1 * 1024 * 1024, // 1MB
		AllowedTypes: []string{"image/jpeg", "image/png"},
		Required:     false,
	}
)

// File size constants for easy reference
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// CreateCustomImageValidation creates a custom image validation configuration
func CreateCustomImageValidation(maxSizeMB int, required bool, allowGif bool) validator.FileValidation {
	allowedTypes := []string{"image/jpeg", "image/png"}
	if allowGif {
		allowedTypes = append(allowedTypes, "image/gif")
	}

	return validator.FileValidation{
		MaxSize:      int64(maxSizeMB * MB),
		AllowedTypes: allowedTypes,
		Required:     required,
	}
}

// CreateCustomDocumentValidation creates a custom document validation configuration
func CreateCustomDocumentValidation(maxSizeMB int, required bool, includeImages bool) validator.FileValidation {
	allowedTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	if includeImages {
		allowedTypes = append(allowedTypes, "image/jpeg", "image/png")
	}

	return validator.FileValidation{
		MaxSize:      int64(maxSizeMB * MB),
		AllowedTypes: allowedTypes,
		Required:     required,
	}
}

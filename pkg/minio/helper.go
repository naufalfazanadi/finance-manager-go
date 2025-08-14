package minio

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
)

type GetFullUrlDto struct {
	BucketName      string
	Path            string
	ExpiredInMinute int64
}

// BucketType represents the type of bucket for upload
type BucketType string

const (
	// BucketTypePublic uploads to public bucket (accessible via direct URL)
	BucketTypePublic BucketType = "public"
	// BucketTypePrivate uploads to private bucket (requires signed URL for access)
	BucketTypePrivate BucketType = "private"
)

type UploadPhotoDto struct {
	FileHeader   *multipart.FileHeader
	FolderPrefix string     // e.g., "profile-photo", "gallery", "documents"
	FilePrefix   string     // e.g., "profile_photo", "gallery_image", "document"
	BucketType   BucketType // "public" or "private" - defaults to "public"
}

type UploadPhotoResult struct {
	Path       string
	FullURL    string
	Filename   string
	BucketName string
	IsPrivate  bool
}

type DeletePhotoDto struct {
	PhotoPath  string     // Required: Path to the file to delete
	BucketType BucketType // Optional: "public" or "private" - if empty, will try both buckets
}

// GetFileExtension extracts the file extension from a filename and returns it in lowercase
// This function is used for file operations and fallback extension detection
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// UploadPhotoMinio handles dynamic photo upload to minio with configurable folder, file prefixes, and bucket type
func UploadPhotoMinio(ctx context.Context, params UploadPhotoDto) (*UploadPhotoResult, error) {
	if params.FileHeader == nil {
		return nil, fmt.Errorf("file header cannot be nil")
	}

	// Set default values if not provided
	if params.FolderPrefix == "" {
		params.FolderPrefix = "uploads"
	}
	if params.FilePrefix == "" {
		params.FilePrefix = "file"
	}
	if params.BucketType == "" {
		params.BucketType = BucketTypePublic // Default to public bucket
	}

	// Open the uploaded file
	file, err := params.FileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Get file extension based on content type
	detectedType := http.DetectContentType(fileContent)
	var ext string
	switch detectedType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	default:
		// Fallback to the original filename extension
		ext = GetFileExtension(params.FileHeader.Filename)
		if ext == "" {
			ext = ".jpg" // ultimate fallback
		}
	}

	// Generate unique filename with timestamp
	timestamp := time.Now().Unix()
	now := time.Now()
	folder := fmt.Sprintf("%s/%d/%02d", params.FolderPrefix, now.Year(), now.Month())
	filename := fmt.Sprintf("%s_%d%s", params.FilePrefix, timestamp, ext)

	// Initialize minio client
	minioClient, err := NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	cfg := config.GetConfig()
	var uploadResult ReturnUploadDto
	var bucketName string
	var isPrivate bool

	// Upload to appropriate bucket based on type
	switch params.BucketType {
	case BucketTypePrivate:
		uploadResult, err = minioClient.UploadPrivate(ctx, UploadPrivateDto{
			OriginalName: params.FileHeader.Filename,
			Folder:       folder,
			FileName:     filename,
			File:         fileContent,
		})
		bucketName = cfg.Minio.PrivateBucket
		isPrivate = true
	case BucketTypePublic:
		uploadResult, err = minioClient.UploadPublic(ctx, UploadPublicDto{
			OriginalName: params.FileHeader.Filename,
			Folder:       folder,
			FileName:     filename,
			File:         fileContent,
		})
		bucketName = cfg.Minio.PublicBucket
		isPrivate = false
	default:
		return nil, fmt.Errorf("invalid bucket type: %s", params.BucketType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to upload to minio: %w", err)
	}

	// Generate full URL based on bucket type
	var fullURL string
	if isPrivate {
		// For private files, generate a signed URL with default expiration
		fullURL = minioClient.GetFullUrl(ctx, GetFullUrlDto{
			BucketName:      bucketName,
			Path:            uploadResult.Path,
			ExpiredInMinute: 60, // Default 1 hour expiration for private files
		})
	} else {
		// For public files, generate direct URL using the helper function
		fullURL = GetFullUrl(bucketName, uploadResult.Path)
	}

	return &UploadPhotoResult{
		Path:       uploadResult.Path,
		FullURL:    fullURL,
		Filename:   filename,
		BucketName: bucketName,
		IsPrivate:  isPrivate,
	}, nil
}

// DeletePhotoMinio handles deletion of uploaded photos from minio with bucket type specification
func DeletePhotoMinio(ctx context.Context, params DeletePhotoDto) error {
	if params.PhotoPath == "" {
		return nil // Nothing to delete
	}

	minioClient, err := NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize minio client: %w", err)
	}

	cfg := config.GetConfig()

	// If bucket type is specified, delete from that specific bucket
	if params.BucketType != "" {
		var bucketName string
		switch params.BucketType {
		case BucketTypePublic:
			bucketName = cfg.Minio.PublicBucket
		case BucketTypePrivate:
			bucketName = cfg.Minio.PrivateBucket
		default:
			return fmt.Errorf("invalid bucket type: %s", params.BucketType)
		}

		if err := minioClient.RemoveObjectByPath(ctx, bucketName, params.PhotoPath); err != nil {
			return fmt.Errorf("failed to delete photo from %s bucket: %w", params.BucketType, err)
		}
		return nil
	}

	// If bucket type is not specified, try to delete from both buckets
	// Try public bucket first
	publicErr := minioClient.RemoveObjectByPath(ctx, cfg.Minio.PublicBucket, params.PhotoPath)

	// Try private bucket
	privateErr := minioClient.RemoveObjectByPath(ctx, cfg.Minio.PrivateBucket, params.PhotoPath)

	// If both failed, return an error
	if publicErr != nil && privateErr != nil {
		return fmt.Errorf("failed to delete photo from both buckets - public: %v, private: %v", publicErr, privateErr)
	}

	// At least one deletion succeeded
	return nil
}

// GetSignedURL generates a signed URL for private files with custom expiration
func GetSignedURL(ctx context.Context, photoPath string, expiredInMinutes int64) (string, error) {
	if photoPath == "" {
		return "", fmt.Errorf("photo path cannot be empty")
	}

	minioClient, err := NewClient()
	if err != nil {
		return "", fmt.Errorf("failed to initialize minio client: %w", err)
	}

	cfg := config.GetConfig()
	return minioClient.GetFullUrl(ctx, GetFullUrlDto{
		BucketName:      cfg.Minio.PrivateBucket,
		Path:            photoPath,
		ExpiredInMinute: expiredInMinutes,
	}), nil
}

// GetFullUrl generates a URL for accessing files in minio buckets
// For private buckets: generates signed URLs with expiration
// For public buckets: generates direct concatenated URLs
//
// Parameters:
//   - bucketName: The name of the minio bucket
//   - path: The file path within the bucket
//   - expiredInMinute: Optional expiration time in minutes (only used for private buckets)
//
// Usage examples:
//   - Public bucket: GetFullUrl("public-bucket", "images/photo.jpg")
//   - Private bucket: GetFullUrl("private-bucket", "docs/secret.pdf")
//   - Private bucket with custom expiration: GetFullUrl("private-bucket", "docs/secret.pdf", 120)
func GetFullUrl(bucketName, path string, expiredInMinute ...int64) string {
	cfg := config.GetConfig()

	// Check if this is a private bucket
	if bucketName == cfg.Minio.PrivateBucket {
		// For private bucket, we need to use the minio client to generate signed URL
		ctx := context.Background()
		minioClient, err := NewClient()
		if err != nil {
			// Fallback to direct URL if client creation fails
			return fmt.Sprintf("https://%s/%s/%s", cfg.Minio.Endpoint, bucketName, path)
		}

		// Set default expiration or use provided value
		expiration := int64(60) // Default 1 hour
		if len(expiredInMinute) > 0 && expiredInMinute[0] > 0 {
			expiration = expiredInMinute[0]
		}

		return minioClient.GetFullUrl(ctx, GetFullUrlDto{
			BucketName:      bucketName,
			Path:            path,
			ExpiredInMinute: expiration,
		})
	}

	// For public bucket or any other bucket, use direct URL concatenation
	return fmt.Sprintf("https://%s/%s/%s", cfg.Minio.Endpoint, bucketName, path)
}

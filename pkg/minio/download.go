package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type DownloadObject struct {
	BucketName      string
	ObjectName      string
	FilePath        string
	ExpiredInMinute int64
}

type GetPrivateDto struct {
	FilePath        string
	ExpiredInMinute int64
}

type ReturnGetPrivateDto struct {
	Path        string
	FullUrl     string
	BucketName  string
	ContentType string
	FileSize    int64
	FileName    string
}

func (c *client) DownloadObject(ctx context.Context, object DownloadObject) error {
	return c.minioClient.FGetObject(ctx, object.BucketName, object.ObjectName, object.FilePath, minio.GetObjectOptions{})
}

func (c *client) GetObject(ctx context.Context, object DownloadObject) (result *minio.Object, err error) {
	return c.minioClient.GetObject(ctx, object.BucketName, object.ObjectName, minio.GetObjectOptions{})
}

func (c *client) RemoveObject(ctx context.Context, object DownloadObject) (err error) {
	return c.minioClient.RemoveObject(ctx, object.BucketName, object.ObjectName, minio.RemoveObjectOptions{})
}

func (c *client) RemoveObjectByPath(ctx context.Context, bucket string, objectPath string) (err error) {
	// For profile photos, we use the public bucket by default
	// The objectPath contains the relative path within the bucket
	return c.minioClient.RemoveObject(ctx, bucket, objectPath, minio.RemoveObjectOptions{})
}

func (c *client) GetObjectPrivate(ctx context.Context, object DownloadObject) (url string, err error) {
	expiry := time.Minute * time.Duration(5) // 5 min.

	if object.ExpiredInMinute != 0 && object.ExpiredInMinute > 0 {
		expiry = time.Minute * time.Duration(object.ExpiredInMinute)
	}

	presignedURL, err := c.minioClient.PresignedGetObject(ctx, privateBucket, object.ObjectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (c *client) GetPrivate(ctx context.Context, object GetPrivateDto) (result ReturnGetPrivateDto, err error) {
	expiry := time.Minute * time.Duration(5) // default 5 min.

	if object.ExpiredInMinute != 0 && object.ExpiredInMinute > 0 {
		expiry = time.Minute * time.Duration(object.ExpiredInMinute)
	}

	presignedURL, err := c.minioClient.PresignedGetObject(ctx, privateBucket, object.FilePath, expiry, nil)

	if err != nil {
		return ReturnGetPrivateDto{}, err
	}

	res := ReturnGetPrivateDto{
		Path:    object.FilePath,
		FullUrl: presignedURL.String(),
	}

	return res, nil
}

func (c *client) GetFullUrl(ctx context.Context, params GetFullUrlDto) (fullUrl string) {
	// Check if this is a private bucket
	if params.BucketName == privateBucket {
		// For private bucket, generate signed URL
		expiry := time.Minute * time.Duration(5) // default 5 min
		if params.ExpiredInMinute > 0 {
			expiry = time.Minute * time.Duration(params.ExpiredInMinute)
		}

		presignedURL, err := c.minioClient.PresignedGetObject(ctx, params.BucketName, params.Path, expiry, nil)
		if err != nil {
			// Fallback to direct URL if signing fails
			return GetFullUrl(params.BucketName, params.Path)
		}
		return presignedURL.String()
	}

	// For public bucket or any other bucket, use direct URL concatenation
	return GetFullUrl(params.BucketName, params.Path)
}

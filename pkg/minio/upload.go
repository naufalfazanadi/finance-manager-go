package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
)

type UploadObject struct {
	BucketName   string
	OriginalName string
	ObjectName   string
	FilePath     string
	ContentType  string
	Reader       io.Reader
	Size         int64
}

type UploadPublicDto struct {
	OriginalName string
	Folder       string
	FileName     string
	File         []byte
}

type UploadPublicFromUrlDto struct {
	OriginalName string
	Folder       string
	FileName     string
	Url          string
}

type UploadPrivateDto struct {
	OriginalName string
	Folder       string
	FileName     string
	File         []byte
}

type ReturnUploadDto struct {
	Path        string `json:"path"`
	FullUrl     string `json:"full_url"`
	BucketName  string `json:"bucket_name"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
	FileName    string `json:"file_name"`
}

func (c *client) UploadObject(ctx context.Context, object UploadObject) (result string, err error) {
	info, err := c.minioClient.PutObject(ctx, object.BucketName, object.ObjectName, object.Reader, object.Size, minio.PutObjectOptions{ContentType: object.ContentType})
	if err != nil {
		return "", err
	}

	result = fmt.Sprintf("%s/%s/%s", c.minioClient.EndpointURL().String(), object.BucketName, info.Key)
	return result, nil
}

func (c *client) UploadObjects(ctx context.Context, objects []UploadObject) (result []string, err error) {
	for _, object := range objects {
		info, err := c.minioClient.PutObject(ctx, object.BucketName, object.ObjectName, object.Reader, object.Size, minio.PutObjectOptions{ContentType: object.ContentType})
		if err != nil {
			return nil, err
		}

		result = append(result, fmt.Sprintf("%s/%s/%s", c.minioClient.EndpointURL().String(), object.BucketName, info.Key))
	}

	return result, nil
}

func (c *client) Upload(ctx context.Context, object UploadObject) (result string, err error) {
	info, err := c.minioClient.PutObject(ctx, object.BucketName, object.ObjectName, object.Reader, object.Size, minio.PutObjectOptions{ContentType: object.ContentType})
	if err != nil {
		return "", err
	}
	result = fmt.Sprintf("%s/%s", object.BucketName, info.Key)
	return result, nil
}

func (c *client) Uploads(ctx context.Context, objects []UploadObject) (result []string, err error) {
	for _, object := range objects {
		info, err := c.minioClient.PutObject(ctx, object.BucketName, object.ObjectName, object.Reader, object.Size, minio.PutObjectOptions{ContentType: object.ContentType})
		if err != nil {
			return nil, err
		}

		result = append(result, fmt.Sprintf("%s/%s", object.BucketName, info.Key))
	}

	return result, nil
}

func (c *client) UploadPrivate(ctx context.Context, object UploadPrivateDto) (result ReturnUploadDto, err error) {
	contentType := http.DetectContentType(object.File)
	fileSize := int64(len(object.File))
	file := bytes.NewReader(object.File)
	objectName := fmt.Sprintf("%s/", directory)
	if object.Folder != "" {
		objectName += fmt.Sprintf("%s/", object.Folder)
	}
	objectName += fmt.Sprintf("%s", strings.ToLower(object.FileName))
	objectName = strings.ReplaceAll(objectName, " ", "_")

	info, err := c.minioClient.PutObject(ctx, privateBucket, objectName, file, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return ReturnUploadDto{}, err
	}

	downloadObject := GetPrivateDto{
		FilePath:        info.Key,
		ExpiredInMinute: 3,
	}
	getPrivateFile, _ := c.GetPrivate(ctx, downloadObject)

	result = ReturnUploadDto{
		Path:        fmt.Sprintf("%s", info.Key),
		FullUrl:     getPrivateFile.FullUrl,
		BucketName:  privateBucket,
		FileSize:    fileSize,
		ContentType: contentType,
		FileName:    object.FileName,
	}

	return result, nil
}

func (c *client) UploadPublic(ctx context.Context, object UploadPublicDto) (result ReturnUploadDto, err error) {
	contentType := http.DetectContentType(object.File)
	fileSize := int64(len(object.File))
	file := bytes.NewReader(object.File)
	objectName := fmt.Sprintf("%s/", directory)
	if object.Folder != "" {
		objectName += fmt.Sprintf("%s/", object.Folder)
	}
	objectName += fmt.Sprintf("%s", strings.ToLower(object.FileName))
	objectName = strings.ReplaceAll(objectName, " ", "_")

	info, err := c.minioClient.PutObject(ctx, publicBucket, objectName, file, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return ReturnUploadDto{}, err
	}

	result = ReturnUploadDto{
		Path:        fmt.Sprintf("%s", info.Key),
		FullUrl:     fmt.Sprintf("%s/%s/%s", c.minioClient.EndpointURL().String(), publicBucket, info.Key),
		BucketName:  publicBucket,
		FileSize:    fileSize,
		ContentType: contentType,
		FileName:    object.FileName,
	}

	return result, nil
}

func (c *client) UploadPublicFromUrl(ctx context.Context, object UploadPublicFromUrlDto) (result ReturnUploadDto, err error) {
	// Note: This function requires implementation of URL to byte conversion
	// For now, returning an error to indicate it's not implemented
	return ReturnUploadDto{}, fmt.Errorf("UploadPublicFromUrl not implemented yet")
}

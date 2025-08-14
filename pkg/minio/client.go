package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
)

var (
	endPoint      string
	accessKey     string
	secretKey     string
	useSSL        bool
	privateBucket string
	publicBucket  string
	directory     string
)

var storeIoClient Client

type Client interface {
	UploadObject(ctx context.Context, object UploadObject) (result string, err error)
	DownloadObject(ctx context.Context, object DownloadObject) error
	UploadObjects(ctx context.Context, objects []UploadObject) (result []string, err error)
	GetObject(ctx context.Context, object DownloadObject) (result *minio.Object, err error)
	GetObjectPrivate(ctx context.Context, object DownloadObject) (result string, err error)
	RemoveObject(ctx context.Context, object DownloadObject) (err error)
	RemoveObjectByPath(ctx context.Context, bucket string, objectPath string) (err error)
	Upload(ctx context.Context, object UploadObject) (result string, err error)
	Uploads(ctx context.Context, objects []UploadObject) (result []string, err error)
	// New Implementation
	GetFullUrl(ctx context.Context, params GetFullUrlDto) (fullUrl string)
	UploadPrivate(ctx context.Context, object UploadPrivateDto) (result ReturnUploadDto, err error)
	UploadPublic(ctx context.Context, object UploadPublicDto) (result ReturnUploadDto, err error)
	UploadPublicFromUrl(ctx context.Context, object UploadPublicFromUrlDto) (result ReturnUploadDto, err error)
	GetPrivate(ctx context.Context, object GetPrivateDto) (result ReturnGetPrivateDto, err error)
}

type client struct {
	minioClient *minio.Client
}

func NewClient() (Client, error) {
	if storeIoClient != nil {
		return storeIoClient, nil
	}

	cfg := config.GetConfig()
	endPoint = cfg.Minio.Endpoint
	accessKey = cfg.Minio.AccessKey
	secretKey = cfg.Minio.SecretKey
	useSSL = cfg.Minio.UseSSL
	privateBucket = cfg.Minio.PrivateBucket
	publicBucket = cfg.Minio.PublicBucket
	directory = cfg.Minio.Directory

	cl, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Initialize policy condition config.
	policy := minio.NewPostPolicy()

	// Apply upload policy restrictions:
	policy.SetBucket(privateBucket)
	policy.SetExpires(time.Now().UTC().Add(2 * 60)) // expires in 2 minutes

	storeIoClient = &client{minioClient: cl}
	return storeIoClient, nil
}

func GetStoreIoClient() Client {
	return storeIoClient
}

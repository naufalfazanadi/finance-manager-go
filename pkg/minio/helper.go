package minio

import (
	"context"
	"fmt"

	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
)

type GetFullUrlDto struct {
	BucketName      string
	Path            string
	ExpiredInMinute int64
}

func (c *client) GetFullUrl(ctx context.Context, params GetFullUrlDto) (fullUrl string) {
	// default full url kalo bucket nya ga ke define private atau public
	fullUrl = params.Path
	expired := params.ExpiredInMinute
	if expired == 0 {
		expired = 3
	}

	cfg := config.GetConfig()
	if params.BucketName == cfg.Minio.PrivateBucket {
		downloadObject := GetPrivateDto{
			FilePath:        params.Path,
			ExpiredInMinute: expired,
		}
		getPrivateFile, _ := c.GetPrivate(ctx, downloadObject)

		fullUrl = getPrivateFile.FullUrl

		return
	} else if params.BucketName == cfg.Minio.PublicBucket {
		fullUrl = fmt.Sprintf("https://%s/%s/%s", cfg.Minio.Endpoint, params.BucketName, params.Path)

		return
	}

	return
}

func GetFullUrl(bucketName, path string) string {
	cfg := config.GetConfig()
	return fmt.Sprintf("https://%s/%s/%s", cfg.Minio.Endpoint, bucketName, path)
}

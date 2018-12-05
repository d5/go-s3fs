package s3fs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Client struct {
	s3Svc      *s3.S3
	bucket     string
	options    *ClientOptions
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// NewClient returns a new file store client.
func NewClient(bucket string, clientOptions *ClientOptions, awsSession client.ConfigProvider, awsConfigs ...*aws.Config) *Client {
	s3Svc := s3.New(awsSession, awsConfigs...)

	uploader := s3manager.NewUploaderWithClient(s3Svc, func(u *s3manager.Uploader) {
		if clientOptions == nil || clientOptions.UploadOptions == nil {
			return
		}
		if clientOptions.UploadOptions.PartSize != nil {
			u.PartSize = *clientOptions.UploadOptions.PartSize
		}
		if clientOptions.UploadOptions.Concurrency != nil {
			u.Concurrency = *clientOptions.UploadOptions.Concurrency
		}
		if clientOptions.UploadOptions.LeavePartsOnError != nil {
			u.LeavePartsOnError = *clientOptions.UploadOptions.LeavePartsOnError
		}
		if clientOptions.UploadOptions.MaxUploadParts != nil {
			u.MaxUploadParts = *clientOptions.UploadOptions.MaxUploadParts
		}
	})

	downloader := s3manager.NewDownloaderWithClient(s3Svc, func(d *s3manager.Downloader) {
		if clientOptions == nil || clientOptions.DownloadOptions == nil {
			return
		}
		if clientOptions.DownloadOptions.PartSize != nil {
			d.PartSize = *clientOptions.DownloadOptions.PartSize
		}
		if clientOptions.DownloadOptions.Concurrency != nil {
			d.Concurrency = *clientOptions.DownloadOptions.Concurrency
		}
	})

	return &Client{
		s3Svc:      s3Svc,
		bucket:     bucket,
		options:    clientOptions,
		uploader:   uploader,
		downloader: downloader,
	}
}

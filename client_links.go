package s3fs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GenerateURLWithExpire creates a direct link (URL) for the file.
// Anyone with this link will be able to access the file (through HTTP GET) regardless of S3 bucket permission settings.
// The link created with this function will return a permission error after it expires.
func (c *Client) GenerateURLWithExpire(path string, expire time.Duration) (string, error) {
	if expire <= 0 {
		return "", fmt.Errorf("expire is zero or negative")
	}

	req, _ := c.s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
	})

	urlStr, err := req.Presign(expire)
	if err != nil {
		return "", fmt.Errorf("error presigning request (%s): %s", path, err.Error())
	}

	return urlStr, nil
}

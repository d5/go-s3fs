package s3fs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Delete deletes a file in path.
// Because the way S3 handles deletion (adding delete-marker for versioning),
// this function does not return any errors when the file does not exist.
func (c *Client) Delete(path string) error {
	_, err := c.s3Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("error deleting file (%s): %s", path, err.Error())
	}

	return nil
}

package s3fs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DeleteAll deletes all files that have pathPrefix as path prefix.
func (c *Client) DeleteAll(pathPrefix string) error {
	deleter := s3manager.NewBatchDeleteWithClient(c.s3Svc)
	deleteIt := s3manager.NewDeleteListIterator(c.s3Svc, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(pathPrefix),
	})

	err := deleter.Delete(aws.BackgroundContext(), deleteIt)
	if err != nil {
		return fmt.Errorf("error deleting in batch (%s): %s", pathPrefix, err.Error())
	}

	return nil
}

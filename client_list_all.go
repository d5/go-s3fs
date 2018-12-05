package s3fs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ListAll returns all files and directories that has pathPrefix as path prefix.
func (c *Client) ListAll(pathPrefix string) ([]ListItem, error) {
	listItems := make([]ListItem, 0)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(pathPrefix),
	}
	err := c.s3Svc.ListObjectsV2Pages(input, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, entry := range output.Contents {
			if entry == nil || entry.Key == nil || *entry.Key == "" {
				continue
			}

			listItems = append(listItems, ListItem{
				Path:  *entry.Key,
				IsDir: (*entry.Key)[len(*entry.Key)-1] == '/',
			})
		}

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating list pages (%s): %s", pathPrefix, err.Error())
	}

	return listItems, nil
}

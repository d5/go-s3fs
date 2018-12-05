package s3fs

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (c *Client) Write(file *File) error {
	if file == nil {
		return errors.New("file is nil")
	}

	// zero-sized file is allowed
	if file.Data == nil {
		file.Data = make([]byte, 0)
	}

	input := &s3manager.UploadInput{
		Bucket:               aws.String(c.bucket),
		Key:                  aws.String(file.Path),
		Body:                 bytes.NewReader(file.Data),
		ContentType:          file.ContentType,
		ACL:                  file.ACL,
		ServerSideEncryption: file.Encryption,
		StorageClass:         file.StorageClass,
	}

	// encode metadata
	encodedMetadata, err := encodeFileMetadata(file.Metadata)
	if err != nil {
		return fmt.Errorf("error encoding file metadata (%s): %s", file.Path, err.Error())
	}
	input.Metadata = encodedMetadata

	if _, err := c.uploader.Upload(input); err != nil {
		return fmt.Errorf("error uploading file (%s): %s", file.Path, err.Error())
	}

	return nil
}

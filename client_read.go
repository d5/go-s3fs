package s3fs

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Read retrieves a file data and its meta data.
// It returns ErrFileNotFound error if there's no file in path.
func (c *Client) Read(path string) (*File, error) {
	// retrieve metadata
	// TODO: maybe using downloader for non-large files is overkill (or even less efficient).
	// One possible approach is to set the threshold (e.g. 5MB) and use that in Range (e.g. "bytes=0-5242880").
	// Then, using response ContentRange (e.g. "bytes 0-5242880/10485760"),
	// we can determine whether we need to download the rest of the file using Downloader or not.
	res, err := c.s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
		Range:  aws.String("bytes=0-0"),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return nil, ErrFileNotFound
			}
			if reqFail, ok := awsErr.(awserr.RequestFailure); ok && reqFail.StatusCode() == http.StatusRequestedRangeNotSatisfiable {
				// zero-sized file; re-try with no range
				res, err = c.s3Svc.GetObject(&s3.GetObjectInput{
					Bucket: aws.String(c.bucket),
					Key:    aws.String(path),
				})
				if err != nil {
					return nil, fmt.Errorf("error retrieving file metadata (%s): %s", path, err.Error())
				}

				return makeFile(path, res, []byte{})
			}
		}

		return nil, fmt.Errorf("error retrieving file metadata (%s): %s", path, err.Error())
	}

	// download data
	buf := &aws.WriteAtBuffer{}
	_, err = c.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return nil, ErrFileNotFound
		}

		return nil, fmt.Errorf("error downloading file (%s): %s", path, err.Error())
	}

	return makeFile(path, res, buf.Bytes())
}

func makeFile(path string, output *s3.GetObjectOutput, data []byte) (*File, error) {
	file := &File{
		Path:         path,
		Data:         data,
		ACL:          nil, // TODO: how do I get this?
		ContentType:  output.ContentType,
		Encryption:   output.ServerSideEncryption,
		StorageClass: output.StorageClass,
	}

	decodedMetadata, err := decodeFileMetadata(output.Metadata)
	if err != nil {
		return nil, fmt.Errorf("error decoding file metadata (%s): %s", path, err.Error())
	}
	file.Metadata = decodedMetadata

	return file, nil
}

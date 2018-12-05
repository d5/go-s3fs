package s3fs

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

var (
	fileMetadataKeyRE      = regexp.MustCompile(`^[\w\-]+$`)
	valueEncodingPrefix    = "+B64E|"
	valueEncodingPrefixLen = len(valueEncodingPrefix)
)

// File represents a file stored in S3 bucket.
type File struct {
	// File path in the storage.
	Path string

	// File data.
	Data []byte

	// A map of metadata to store with the file in S3.
	// This is optional.
	Metadata map[string]string

	// Standard MIME type describing the format of the file data.
	// This is optional, and, will be set to "binary/octet-stream" by default.
	ContentType *string

	// Canned ACL to apply to the file.
	// Possible values include
	//  - "private"
	//  - "public-read"
	//  - "public-read-write"
	//  - "authenticated-read"
	//  - "aws-exec-read"
	//  - "bucket-owner-read"
	//  - "bucket-owner-full-control"
	// This is optional.
	ACL *string

	// Server-side encryption algorithm used when storing file.
	// Possible values include "AES256", "aws:kms".
	// This is optional.
	Encryption *string

	// Type of storage to use for the file.
	// Possible values include "STANDARD", "REDUCED_REDUNDANCY", "STANDARD_IA".
	// This is optional, and, will be set to "STANDARD" by default.
	StorageClass *string
}

func encodeFileMetadata(input map[string]string) (map[string]*string, error) {
	var output map[string]*string = nil

	if input != nil && len(input) > 0 {
		output = make(map[string]*string)

		for key, value := range input {
			// allowed characters in key: [\w\-]+
			if !fileMetadataKeyRE.MatchString(key) {
				return nil, fmt.Errorf("invalid character in key: %q", key)
			}
			keyLower := strings.ToLower(key)

			// Value
			if !isASCII(value) {
				// value encoded needed
				output[keyLower] = aws.String(valueEncodingPrefix + base64.URLEncoding.EncodeToString([]byte(value)))
			} else {
				output[keyLower] = aws.String(value)
			}
		}
	}

	return output, nil
}

func decodeFileMetadata(input map[string]*string) (map[string]string, error) {
	var output map[string]string = nil

	if input != nil && len(input) > 0 {
		output = make(map[string]string)

		for key, value := range input {
			keyLower := strings.ToLower(key)

			if value == nil {
				output[keyLower] = ""
			} else {
				if strings.HasPrefix(*value, valueEncodingPrefix) {
					buf, err := base64.URLEncoding.DecodeString((*value)[valueEncodingPrefixLen:])
					if err != nil {
						return nil, fmt.Errorf("error decoding value (key=%s): %s", key, err.Error())
					}
					output[keyLower] = string(buf)
				} else {
					output[keyLower] = *value
				}
			}

		}
	}

	return output, nil
}

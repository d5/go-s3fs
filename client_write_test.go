package s3fs_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/d5/go-s3fs"
	"github.com/stretchr/testify/assert"
)

func TestClient_Write(t *testing.T) {
	client := MustCreateS3Client()

	// nil error
	err := client.Write(nil)
	assert.NotNil(t, err)

	// file.path == ""
	err = client.Write(&s3fs.File{
		Data: []byte("hello"),
	})
	assert.NotNil(t, err)

	// test files
	testFilePath1 := randomTestFilePath()
	defer func() {
		_ = client.Delete(testFilePath1)
	}()

	// simple
	file := &s3fs.File{
		Path: testFilePath1,
		Data: []byte("test foo data"),
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, file)

	// overwrite
	file = &s3fs.File{
		Path: testFilePath1,
		Data: []byte("test bar data"),
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, file)

	// zero-sized data is allowed
	file = &s3fs.File{
		Path: testFilePath1,
	}
	err = client.Write(file)
	assertEqualFile(t, client, file)
	file = &s3fs.File{
		Path: testFilePath1,
		Data: make([]byte, 0),
	}
	err = client.Write(file)
	assertEqualFile(t, client, file)

	// metadata
	file = &s3fs.File{
		Path: testFilePath1,
		Data: []byte("bar data"),
		Metadata: map[string]string{
			"key1_bar": "value-more",
			"key2-bar": "VALUE2 and MORE",
			"key3-bar": "こんにちは世界",
		},
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, file)

	// invalid metadata key
	file = &s3fs.File{
		Path: testFilePath1,
		Data: []byte("bar data"),
		Metadata: map[string]string{
			"key1 bar": "value", // space
		},
	}
	err = client.Write(file)
	assert.NotNil(t, err)
	file = &s3fs.File{
		Path: testFilePath1,
		Data: []byte("bar data"),
		Metadata: map[string]string{
			"키1": "value", // non-ASCII
		},
	}
	err = client.Write(file)
	assert.NotNil(t, err)

	// metadata key stored in lowercase
	file = &s3fs.File{
		Path: testFilePath1,
		Data: []byte("bar data"),
		Metadata: map[string]string{
			"Key1-Part1": "value1",
			"Key2-Part2": "value2 MORE",
		},
	}
	fileMKLower := &s3fs.File{
		Path: testFilePath1,
		Data: []byte("bar data"),
		Metadata: map[string]string{
			"key1-part1": "value1",
			"key2-part2": "value2 MORE",
		},
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, fileMKLower)

	// ACL
	file = &s3fs.File{
		Path: testFilePath1,
		Data: []byte("test bar data"),
		ACL:  s3fs.StringPtr(s3.ObjectCannedACLPublicRead),
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, file)

	// ContentType
	file = &s3fs.File{
		Path:        testFilePath1,
		Data:        []byte(`{"foo": "bar"}"`),
		ContentType: s3fs.StringPtr("application/json"),
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, file)

	// Encryption
	file = &s3fs.File{
		Path:       testFilePath1,
		Data:       []byte(`{"foo": "bar"}"`),
		Encryption: s3fs.StringPtr(s3fs.S3EncryptionAES256),
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, file)
	// wrong encryption type
	file = &s3fs.File{
		Path:       testFilePath1,
		Data:       []byte(`{"foo": "bar"}"`),
		Encryption: s3fs.StringPtr(s3fs.S3EncryptionAES256 + "_unsupported"),
	}
	err = client.Write(file)
	assert.NotNil(t, err)

	// Storage class
	file = &s3fs.File{
		Path:         testFilePath1,
		Data:         []byte(`{"foo": "bar"}"`),
		StorageClass: s3fs.StringPtr(s3fs.S3StorageClassReducedRedundancy),
	}
	err = client.Write(file)
	assert.Nil(t, err)
	assertEqualFile(t, client, file)
}

func assertEqualFile(t *testing.T, client *s3fs.Client, expected *s3fs.File) {
	actual, err := client.Read(expected.Path)
	if assert.Nil(t, err) && assert.NotNil(t, actual) {
		// TODO: Client.Read() does not retrieve "canned" ACL property in the current implementation
		expected.ACL = nil

		// ContentType is set to "binary/octet-stream" by default
		if expected.ContentType == nil {
			expected.ContentType = s3fs.StringPtr("binary/octet-stream")
		}

		assert.Equal(t, expected, actual)
	}
}

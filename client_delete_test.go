package s3fs_test

import (
	"testing"

	"github.com/d5/go-s3fs"
	"github.com/stretchr/testify/assert"
)

func TestClient_Delete(t *testing.T) {
	client := MustCreateS3Client()
	testFilePath := randomTestFilePath()

	// empty path
	err := client.Delete("")
	assert.NotNil(t, err)

	// delete non existing file: no error
	err = client.Delete(testFilePath)
	assert.Nil(t, err)

	// delete actual normal file
	testFile := &s3fs.File{Path: testFilePath, Data: []byte("hello")}
	_ = client.Write(testFile)
	assertEqualFile(t, client, testFile) // make sure file was created
	err = client.Delete(testFilePath)    // delete it
	assert.Nil(t, err)
	_, err = client.Read(testFilePath)
	assert.Equal(t, s3fs.ErrFileNotFound, err) // file should be deleted
}

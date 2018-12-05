package s3fs_test

import (
	"testing"

	"github.com/d5/go-s3fs"
	"github.com/stretchr/testify/assert"
)

func TestClient_Read(t *testing.T) {
	client := MustCreateS3Client()

	testFilePath1 := randomTestFilePath()

	// error: empty path
	_, err := client.Read("")
	assert.NotNil(t, err)

	// read non-existing file
	_, err = client.Read(testFilePath1)
	assert.Equal(t, s3fs.ErrFileNotFound, err)

	// NOTE: most Read functions are tested through TestClient_Write()
}

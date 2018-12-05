package s3fs_test

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/d5/go-httputil"
	"github.com/d5/go-s3fs"
	"github.com/stretchr/testify/assert"
)

func TestClient_GenerateURLWithExpire(t *testing.T) {
	client := MustCreateS3Client()

	// create test file
	testFile := &s3fs.File{
		Path:        randomTestFilePath(),
		Data:        []byte("hello world!"),
		ContentType: aws.String("text/plain"),
	}
	err := client.Write(testFile)
	if assert.Nil(t, err) {
		defer client.Delete(testFile.Path)
	}

	// error: expire <= 0
	_, err = client.GenerateURLWithExpire(testFile.Path, 0)
	assert.NotNil(t, err)
	_, err = client.GenerateURLWithExpire(testFile.Path, -1*time.Second)
	assert.NotNil(t, err)

	// simple link (expire = 10m)
	link, err := client.GenerateURLWithExpire(testFile.Path, 1*time.Minute)
	assert.Nil(t, err)
	testHTTPGet(t, link, http.StatusOK, testFile.ContentType, testFile.Data)

	// expired link (expire = 1ms)
	link, err = client.GenerateURLWithExpire(testFile.Path, 1*time.Millisecond)
	assert.Nil(t, err)
	time.Sleep(10 * time.Millisecond) // pause 10ms
	testHTTPGet(t, link, http.StatusForbidden, nil, nil)
}

func testHTTPGet(t *testing.T, link string, expectedStatus int, expectedContentType *string, expectedBody []byte) {
	httpRes, err := httputil.NewClient().Get(link)
	if !assert.Nil(t, err) {
		return
	}
	defer httpRes.Body.Close()

	assert.Equal(t, expectedStatus, httpRes.StatusCode)

	if expectedContentType != nil {
		assert.Equal(t, *expectedContentType, httpRes.Header.Get("Content-Type"))
	}

	if expectedBody != nil {
		resBody, err := ioutil.ReadAll(httpRes.Body)
		assert.Nil(t, err)
		assert.Equal(t, expectedBody, resBody)
	}
}

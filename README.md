# s3fs

[![GoDoc](https://godoc.org/github.com/d5/go-s3fs?status.svg)](https://godoc.org/github.com/d5/go-s3fs)


A simple S3-based file store implementation for Go.

Example:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/d5/go-s3fs"
)

func main() {
	// create client
	client := s3fs.NewClient(
		os.Getenv("S3_BUCKET"), // S3 bucket name
		nil,                  // use default client options
		s3fs.NewAWSSession(), // AWS session
		s3fs.NewAWSConfigWithStaticCredentialsAndRegion(
			os.Getenv("AWS_ACCESS_KEY_ID"),     // AWS access key
			os.Getenv("AWS_SECRET_ACCESS_KEY"), // AWS secret key
			os.Getenv("S3_BUCKET_REGION")))     // S3 bucket region

	// write a file
	err := client.Write(&s3fs.File{
		Path:        "test/foo/bar",
		Data:        []byte("hello world!"),
		ContentType: s3fs.StringPtr("text/plain"),
	})
	if err != nil {
		panic(err)
	}

	// read a file
	file, err := client.Read("test/foo/bar")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(file.Data))

	// get a pre-signed URL for the file
	link, err := client.GenerateURLWithExpire("test/foo/bar", 10*time.Minute)
	if err != nil {
		panic(err)
	}
	fmt.Println(link)
}
```

To run tests, you need to create your own S3 bucket and set the following environment variables.

```bash
export AWS_ACCESS_KEY_ID="<your AWS access key>"
export AWS_SECRET_ACCESS_KEY="<your AWS secret key>"
export TEST_S3_BUCKET="<your S3 bucket name>"
export TEST_S3_BUCKET_REGION="<your S3 bucket region>"

# run tests
make test
```

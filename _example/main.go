package main

import (
	"fmt"
	"os"
	"time"

	"github.com/d5/go-s3fs"
)

func main() {
	// create FS client
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

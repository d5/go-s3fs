package s3fs_test

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/d5/go-s3fs"
	"github.com/satori/go.uuid"
)

var (
	awsAccessKey       = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey       = os.Getenv("AWS_SECRET_ACCESS_KEY")
	testS3Bucket       = os.Getenv("TEST_S3_BUCKET")
	testS3BucketRegion = os.Getenv("TEST_S3_BUCKET_REGION")
)

func MustCreateS3Client() *s3fs.Client {
	if awsAccessKey == "" || awsSecretKey == "" || testS3Bucket == "" || testS3BucketRegion == "" {
		panic(fmt.Errorf("env vars for test not set: $AWS_ACCESS_KEY_ID, $AWS_SECRET_ACCESS_KEY, $TEST_S3_BUCKET, $TEST_S3_BUCKET_REGION"))
	}

	awsConfig := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")).
		WithRegion(testS3BucketRegion)

	return s3fs.NewClient(testS3Bucket, nil, session.Must(session.NewSession()), awsConfig)
}

func remoteTestDir() string {
	return path.Join("s3fs", "test")
}

func randomTestFilePath() string {
	return path.Join(remoteTestDir(), hex.EncodeToString(uuid.NewV4().Bytes()))
}

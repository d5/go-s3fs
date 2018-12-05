package s3fs

import (
	"unicode"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}

	return true
}

// StringPtr converts string to string pointer.
func StringPtr(v string) *string {
	return &v
}

// NewAWSSession returns a default AWS session.
func NewAWSSession() *session.Session {
	return session.Must(session.NewSession())
}

// NewAWSConfigWithStaticCredentials returns AWS config with static credentials.
func NewAWSConfigWithStaticCredentials(awsAccessKey, awsSecretKey string) *aws.Config {
	return aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""))
}

// NewAWSConfigWithStaticCredentialsAndRegion returns AWS config with static credentials and region.
func NewAWSConfigWithStaticCredentialsAndRegion(awsAccessKey, awsSecretKey, awsRegion string) *aws.Config {
	return NewAWSConfigWithStaticCredentials(awsAccessKey, awsSecretKey).WithRegion(awsRegion)
}

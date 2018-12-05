package s3fs

const (
	S3ACLPrivate                = "private"
	S3ACLPublicRead             = "public-read"
	S3ACLPublicReadWrite        = "public-read-write"
	S3ACLAuthenticatedRead      = "authenticated-read"
	S3ACLAWSExecRead            = "aws-exec-read"
	S3ACLBucketOwnerRead        = "bucket-owner-read"
	S3ACLBucketOwnerFullControl = "bucket-owner-full-control"
)

const (
	S3StorageClassStandard          = "STANDARD"
	S3StorageClassReducedRedundancy = "REDUCED_REDUNDANCY"
	S3StorageClassStandardIA        = "STANDARD_IA"
)

const (
	S3EncryptionAES256 = "AES256"
	S3EncryptionAWSKMS = "aws:kms"
)

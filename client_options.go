package s3fs

type ClientOptions struct {
	UploadOptions   *UploadOptions
	DownloadOptions *DownloadOptions
}

type UploadOptions struct {
	PartSize          *int64
	Concurrency       *int
	LeavePartsOnError *bool
	MaxUploadParts    *int
}

type DownloadOptions struct {
	PartSize    *int64
	Concurrency *int
}

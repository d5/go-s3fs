package s3fs

// ListItem represents a file or directory in ListAll result.
type ListItem struct {
	Path  string // file or directory path (directory path ends with slash)
	IsDir bool   // whether item is directory or not
}

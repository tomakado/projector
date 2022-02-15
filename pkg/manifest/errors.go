package manifest

import "errors"

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrPermissionDenied = errors.New("permission denied")
)

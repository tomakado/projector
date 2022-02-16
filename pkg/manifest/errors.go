package manifest

import "errors"

var (
	// ErrFileNotFound is returned when provider is unable to locate file.
	ErrFileNotFound = errors.New("file not found")

	// ErrPermissionDenied is returned when provider doesn't have permissions to open file.
	ErrPermissionDenied = errors.New("permission denied")
)

package manifest

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type RealFSProvider struct {
	root string
}

func NewRealFSProvider(root string) *RealFSProvider {
	return &RealFSProvider{root: root}
}

func (r *RealFSProvider) Get(filename string) ([]byte, error) {
	fullPath := filepath.Join(r.root, filename)

	f, err := os.OpenFile(fullPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		errToWrap := err

		switch {
		case errors.Is(err, fs.ErrNotExist):
			errToWrap = ErrFileNotFound
		case errors.Is(err, fs.ErrPermission):
			errToWrap = ErrPermissionDenied
		}

		return nil, fmt.Errorf("open %q: %w", fullPath, errToWrap)
	}
	defer f.Close()

	bts, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", fullPath, err)
	}

	return bts, nil
}

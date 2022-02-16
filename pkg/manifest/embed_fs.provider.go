package manifest

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
)

// EmbedFSProvider is an wrapper for embed.FS that provides
// implementation of provider interface accepted by projector.Generator.
type EmbedFSProvider struct {
	fs   *embed.FS
	root string
}

func NewEmbedFSProvider(fs *embed.FS, root string) *EmbedFSProvider {
	return &EmbedFSProvider{
		fs:   fs,
		root: root,
	}
}

func (e *EmbedFSProvider) Get(filename string) ([]byte, error) {
	f, err := e.fs.Open(filepath.Join(e.root, filename))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("open %q: %w", filename, ErrFileNotFound)
		}
		return nil, fmt.Errorf("open %q: %w", filename, err)
	}
	defer f.Close() //nolint:errcheck

	bts, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", filename, err)
	}

	return bts, nil
}

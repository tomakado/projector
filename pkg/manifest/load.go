package manifest

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/tomakado/projector/internal/pkg/verbose"
)

type provider interface {
	Get(filename string) ([]byte, error)
}

func Load(p provider, path string) (*Manifest, error) {
	verbose.Printf("loading manifest %q", path)

	manifestBytes, err := p.Get(path)
	if err != nil {
		return nil, err
	}

	manifest, err := parseManifest(manifestBytes)
	if err != nil {
		return nil, err
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return manifest, nil
}

func parseManifest(src []byte) (*Manifest, error) {
	verbose.Println("parsing manifest")

	var manifest *Manifest
	if err := toml.Unmarshal(src, &manifest); err != nil {
		// TODO wrap custom typed error
		verbose.Println("toml.Unmarshal returned error")
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	return manifest, nil
}

package manifest

import "embed"

// Manifest contains all metadata related to project template and actual steps of project generation.
type Manifest struct {
	fs *embed.FS

	Name    string `toml:"name"`
	Author  string `toml:"author,omitempty"`
	URL     string `toml:"url,omitempty"`
	Version string `toml:"version,omitempty"`

	Steps []Step `toml:"steps"`
}

// EmbeddedFS returns reference to embedded filesystem this manifest belongs to.
func (t *Manifest) EmbeddedFS() *embed.FS {
	return t.fs
}

// WithEmbeddedFS writes fs to manifest and returns reference to manifest.
func (t *Manifest) WithEmbeddedFS(fs *embed.FS) *Manifest {
	t.fs = fs
	return t
}

// Step contains template files to output mapping and/or shell script to execute.
type Step struct {
	Name  string `toml:"name"`
	Files []File `toml:"files"`
	Shell string `toml:"shell"`
}

// File is actually mapping between template file and output file. Also template syntax allowed in Output field.
type File struct {
	Path   string `toml:"path"`
	Output string `toml:"output"`
}

package projector

import "github.com/tomakado/projector/pkg/manifest"

// Config contains all information required to generate project.
type Config struct {
	WorkingDirectory string
	ProjectAuthor    string
	ProjectName      string
	ProjectPackage   string
	Manifest         *manifest.Manifest
	ManifestPath     string
}

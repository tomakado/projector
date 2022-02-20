package projector

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/tomakado/projector/internal/pkg/verbose"
	"github.com/tomakado/projector/pkg/manifest"
)

type CreateConfig struct {
	Config         *Config
	Provider       provider
	PathToManifest string
}

func Create(cfg CreateConfig) error {
	m, err := manifest.Load(cfg.Provider, filepath.Join(cfg.PathToManifest, "projector.toml"))
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	cfg.Config.ManifestPath = cfg.PathToManifest
	cfg.Config.Manifest = m

	if cfg.Config.ProjectPackage == "" {
		cfg.Config.ProjectPackage = cfg.Config.ProjectName

		verbose.Printf(
			"project package name is not provided, using project name as package name (%q)",
			cfg.Config.ProjectPackage,
		)
	}

	if cfg.Config.ProjectAuthor == "" {
		u, err := user.Current()
		if err != nil {
			// TODO wrap custom typed error if possible
			return fmt.Errorf("get current user: %w", err)
		}
		cfg.Config.ProjectAuthor = u.Username

		verbose.Printf(
			"project author is not provided, using current OS user as author (%q)",
			cfg.Config.ProjectAuthor,
		)
	}

	return Generate(cfg.Config, cfg.Provider)
}

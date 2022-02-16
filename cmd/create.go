package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	projector "github.com/tomakado/projector/pkg"
	"github.com/tomakado/projector/pkg/manifest"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create project using specified template",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runCreate,
	}
	cfg            projector.Config
	pathToManifest string
)

func init() {
	createCmd.Flags().StringVarP(&cfg.ProjectName, "name", "n", "my-app", "project name")
	createCmd.Flags().StringVarP(
		&cfg.ProjectPackage,
		"package",
		"p",
		"",
		"project's module name (default same as project name)",
	)
	createCmd.Flags().StringVarP(&cfg.ProjectAuthor, "author", "a", "", "project author (default current OS user)")
	createCmd.Flags().StringVarP(&pathToManifest, "manifest", "m", "", "path to custom template manifest")
}

func runCreate(_ *cobra.Command, args []string) error {
	var p provider
	if pathToManifest != "" {
		p = manifest.NewRealFSProvider(filepath.Dir(pathToManifest))
		cfg.WorkingDirectory = args[0]
	} else {
		p = manifest.NewEmbedFSProvider(&resources, embedRoot)
		cfg.WorkingDirectory = args[1]
	}

	manifest, err := loadManifest(p, pathToManifest)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	cfg.ManifestPath = pathToManifest
	cfg.Manifest = manifest

	if cfg.ProjectPackage == "" {
		cfg.ProjectPackage = cfg.ProjectName
	}

	if cfg.ProjectAuthor == "" {
		u, err := user.Current()
		if err != nil {
			// TODO wrap custom typed error (if possible)
			return fmt.Errorf("get current user: %w", err)
		}

		cfg.ProjectAuthor = u.Name
	}

	return projector.Generate(&cfg, p)
}

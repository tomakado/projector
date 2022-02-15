package cmd

import (
	"embed"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	projector "github.com/tomakado/projector/pkg"
	"github.com/tomakado/projector/pkg/manifest"
)

var (
	createCmd = &cobra.Command{
		Use:  "create",
		Args: cobra.MinimumNArgs(1),
		RunE: createE,
	}
	cfg            projector.Config
	pathToManifest string
)

//go:embed resources/*
var resources embed.FS

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
	createCmd.Flags().StringVarP(&pathToManifest, "file", "f", "", "path to custom template manifest")
}

func createE(_ *cobra.Command, args []string) error {
	if pathToManifest != "" {
		fillConfigForCustomManifest(pathToManifest)
		cfg.WorkingDirectory = args[0]
	} else {
		fillConfigForEmbeddedManifest(args[0], args[1])
		cfg.WorkingDirectory = args[1]
	}

	if cfg.ProjectPackage == "" {
		cfg.ProjectPackage = cfg.ProjectName
	}

	if cfg.ProjectAuthor == "" {
		u, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		cfg.ProjectAuthor = u.Name
	}

	return projector.Generate(&cfg)
}

func fillConfigForCustomManifest(path string) error {
	manifestBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read manifest of custom project template: %w", err)
	}

	manifest, err := parseManifest(manifestBytes)
	if err != nil {
		return err
	}

	cfg.Manifest = manifest
	cfg.ManifestPath = filepath.Dir(path)

	return nil
}

func fillConfigForEmbeddedManifest(templateName, workingDirectory string) error {
	embedManifestPath := fmt.Sprintf("resources/templates/%s/projector.toml", templateName)

	manifestBytes, err := resources.ReadFile(embedManifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest of project template %q: %w", templateName, err)
	}

	manifest, err := parseManifest(manifestBytes)
	if err != nil {
		return err
	}

	cfg.Manifest = manifest.WithEmbeddedFS(&resources)
	cfg.ManifestPath = filepath.Dir(embedManifestPath)

	return nil
}

func parseManifest(src []byte) (*manifest.Manifest, error) {
	var manifest *manifest.Manifest
	if err := toml.Unmarshal(src, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest of project template: %w", err)
	}

	return manifest, nil
}

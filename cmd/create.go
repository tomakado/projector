package cmd

import (
	"embed"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	projector "github.com/tomakado/projector/pkg"
)

var (
	createCmd = &cobra.Command{
		Use:  "create",
		Args: cobra.ExactArgs(2),
		RunE: createE,
	}
	cfg projector.Config
)

//go:embed resources/*
var resources embed.FS

func init() {
	createCmd.Flags().StringVarP(&cfg.ProjectName, "name", "n", "my-app", "project name")
	createCmd.Flags().StringVarP(&cfg.ProjectPackage, "package", "p", "", "project's module name")
	createCmd.Flags().StringVarP(&cfg.ProjectAuthor, "author", "a", "anonymous", "project author") // TODO: use OS username by default!
}

func createE(_ *cobra.Command, args []string) error {
	var (
		templateName    = args[0]
		tplManifestPath = fmt.Sprintf("resources/templates/%s/projector.toml", templateName)
	)

	tplManifestBytes, err := resources.ReadFile(tplManifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest of project template %q: %w", templateName, err)
	}

	var tplManifest *projector.TemplateManifest
	if err := toml.Unmarshal(tplManifestBytes, &tplManifest); err != nil {
		return fmt.Errorf("failed to parse manifest of project template %q: %w", templateName, err)
	}

	cfg.WorkingDirectory = args[1]
	cfg.Manifest = tplManifest.WithEmbeddedFS(&resources)

	if cfg.ProjectPackage == "" {
		cfg.ProjectPackage = cfg.ProjectName
	}

	return projector.Generate(&cfg)
}

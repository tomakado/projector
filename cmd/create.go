package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tomakado/projector/internal/pkg/verbose"
	projector "github.com/tomakado/projector/pkg"
	"github.com/tomakado/projector/pkg/manifest"
)

var (
	createCmd = &cobra.Command{
		Use:   "create [TEMPLATE]",
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
		verbose.Printf("custom manifest filename passed: %q", pathToManifest)
		p = manifest.NewRealFSProvider(filepath.Dir(pathToManifest))
		cfg.WorkingDirectory = args[0]
	} else {
		pathToManifest = args[0]
		verbose.Printf("using manifest name %q in embed fs", pathToManifest)
		p = manifest.NewEmbedFSProvider(&resources, embedRoot)
		cfg.WorkingDirectory = args[1]
	}

	verbose.Printf("working directory = %q", cfg.WorkingDirectory)

	return projector.Create(
		projector.CreateConfig{
			Config:         &cfg,
			Provider:       p,
			PathToManifest: pathToManifest,
		},
	)
}

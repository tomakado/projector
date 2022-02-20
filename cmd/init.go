package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	projector "github.com/tomakado/projector/pkg"
	"github.com/tomakado/projector/pkg/manifest"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create template manifest in current directory (like `create projector` command)",
	RunE:  runInit,
}

func runInit(*cobra.Command, []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	var (
		projectName = path.Base(wd)
		p           = manifest.NewEmbedFSProvider(&resources, embedRoot)
	)

	cfg := &projector.Config{
		ProjectName:      projectName,
		ProjectPackage:   projectName,
		WorkingDirectory: wd,
	}

	return projector.Create(
		projector.CreateConfig{
			Config:         cfg,
			Provider:       p,
			PathToManifest: "projector",
		},
	)
}

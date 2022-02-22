package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomakado/projector/internal/pkg/verbose"
	projector "github.com/tomakado/projector/pkg"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List builtin and cached templates",
	RunE:  runList,
}

func runList(_ *cobra.Command, _ []string) error {
	verbose.Println("traversing templates tree")

	manifests, err := projector.CollectEmbeddedManifests(&resources, embedRoot, ".")
	if err != nil {
		return fmt.Errorf("collect manifests: %w", err)
	}

	for _, m := range manifests {
		fmt.Println(m)
	}

	return nil
}

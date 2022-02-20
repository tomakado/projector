package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomakado/projector/internal/pkg/verbose"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List builtin and cached templates",
	RunE:  runList,
}

func runList(_ *cobra.Command, _ []string) error {
	verbose.Println("traversing templates tree")

	manifests, err := collectManifests(".")
	if err != nil {
		return fmt.Errorf("collect manifests: %w", err)
	}

	for _, m := range manifests {
		fmt.Println(m)
	}

	return nil
}

func collectManifests(root string) ([]string, error) {
	verbose.Printf("reading %q", root)

	var manifests []string

	dirs, err := resources.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read dir %q: %w", root, err)
	}

	for _, entry := range dirs {
		if !entry.IsDir() {
			if entry.Name() == "projector.toml" {
				manifestName := strings.TrimPrefix(root, embedRoot)
				manifests = append(manifests, manifestName) //nolint:staticcheck

				verbose.Printf("projector.toml detected, so registered %q as manifest", manifestName)
			}

			continue
		}

		children, err := collectManifests(path.Join(root, entry.Name()))
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, children...)
	}

	return manifests, nil
}

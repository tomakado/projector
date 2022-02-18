package cmd

import (
	"fmt"
	"path"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tomakado/projector/pkg/manifest"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List builtin and cached templates",
	RunE:  runList,
}

func runList(_ *cobra.Command, _ []string) error {
	manifests, err := collectManifests(".")
	if err != nil {
		return fmt.Errorf("collect manifests: %w", err)
	}

	for _, m := range manifests {
		color.New(color.FgGreen).Printf("%s@%s", m.Name, m.Version)
		fmt.Printf(" by %s\n", m.Author)
		if m.URL != "" {
			color.New(color.FgWhite).Printf("%s\n\n", m.URL)
		} else {
			fmt.Println()
		}
	}

	return nil
}

func collectManifests(root string) ([]manifest.Manifest, error) {
	var manifests []manifest.Manifest

	dirs, err := resources.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read dir %q: %w", root, err)
	}

	p := manifest.NewEmbedFSProvider(&resources, root)

	for _, entry := range dirs {
		if !entry.IsDir() {
			if entry.Name() == "projector.toml" {
				m, err := loadManifest(p, "projector.toml")
				if err != nil {
					return nil, fmt.Errorf("load %q: %w", root, err)
				}

				manifests = append(manifests, *m)
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

package projector

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/tomakado/projector/internal/pkg/verbose"
)

func CollectEmbeddedManifests(fs *embed.FS, embedRoot, root string) ([]string, error) {
	// TODO Make this func universal and use some kind of registry of manifests.
	// Making it after adding support for GitHub as manifest source is the best moment.
	verbose.Printf("reading %q", root)

	var manifests []string

	dirs, err := fs.ReadDir(root)
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

		children, err := CollectEmbeddedManifests(fs, embedRoot, path.Join(root, entry.Name()))
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, children...)
	}

	return manifests, nil
}

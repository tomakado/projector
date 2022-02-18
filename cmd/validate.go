package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tomakado/projector/pkg/manifest"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate [TEMPLATE?]",
		Short: "Validate manifest without performing actions (dry run)",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runValidate,
	}
	manifestNameToValidate string
)

func init() {
	validateCmd.Flags().StringVarP(&manifestNameToValidate, "manifest", "m", "", "manifest name or path to validate")
}

func runValidate(_ *cobra.Command, args []string) error {
	var p provider
	if manifestNameToValidate == "" {
		if len(args) == 0 {
			return fmt.Errorf("builtin manifest name as argument is required if path to manifest is not specified")
		}

		manifestNameToValidate = filepath.Join(args[0], "projector.toml")
		p = manifest.NewEmbedFSProvider(&resources, embedRoot)
	} else {
		p = manifest.NewRealFSProvider(filepath.Dir(manifestNameToValidate))
	}

	_, err := loadManifest(p, manifestNameToValidate)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	fmt.Println("Manifest is valid ✅")

	return nil
}

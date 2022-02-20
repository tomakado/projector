package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tomakado/projector/internal/pkg/verbose"
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
		verbose.Println("manifest name is not passed from flag, trying to read from args")

		if len(args) == 0 {
			return fmt.Errorf("builtin manifest name as argument is required if path to manifest is not specified")
		}

		manifestNameToValidate = filepath.Join(args[0], "projector.toml")
		verbose.Printf("using manifest name %q", manifestNameToValidate)

		p = manifest.NewEmbedFSProvider(&resources, embedRoot)
	} else {
		p = manifest.NewRealFSProvider(filepath.Dir(manifestNameToValidate))
	}

	_, err := manifest.Load(p, manifestNameToValidate)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	fmt.Println("Manifest is valid ✅")

	return nil
}

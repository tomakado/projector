package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tomakado/projector/pkg/manifest"
)

var infoCmd = &cobra.Command{
	Use:   "info [TEMPLATE]",
	Short: "Show meta information about template",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func runInfo(_ *cobra.Command, args []string) error {
	var (
		templateName = args[0]
		p            = manifest.NewEmbedFSProvider(&resources, embedRoot)
	)

	m, err := manifest.Load(p, filepath.Join(templateName, "projector.toml"))
	if err != nil {
		return err
	}

	printManifest(m)

	return nil
}

func printManifest(m *manifest.Manifest) {
	color.New(color.FgGreen).Printf("%s@%s", m.Name, m.Version)
	fmt.Printf(" by %s\n", m.Author)

	if m.URL != "" {
		color.New(color.Bold).Print("URL: ")
		fmt.Println(m.URL)
	} else {
		fmt.Println()
	}

	if m.Description != "" {
		color.New(color.Bold).Print("Description: ")
		fmt.Println(m.Description)
	}
}

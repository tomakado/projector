package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/tomakado/projector/pkg/manifest"
)

const embedRoot = "resources/templates/"

var rootCmd = &cobra.Command{
	Use:   "projector",
	Short: "Projector is a tool for generating project from templates",
	Long: `A flexible, language and framework agnostic tool that allows you to generate projects from templates. 
Projector has some builtin templates, but you can use your custom templates or import third-party templates
from GitHub.`,
	SilenceUsage: true,
}

//go:embed resources/*
var resources embed.FS

type provider interface {
	Get(filename string) ([]byte, error)
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(validateCmd)
}

// Execute runs passed command and handles errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadManifest(p provider, path string) (*manifest.Manifest, error) {
	manifestBytes, err := p.Get(path)
	if err != nil {
		return nil, err
	}

	manifest, err := parseManifest(manifestBytes)
	if err != nil {
		return nil, err
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return manifest, nil
}

func parseManifest(src []byte) (*manifest.Manifest, error) {
	var manifest *manifest.Manifest
	if err := toml.Unmarshal(src, &manifest); err != nil {
		// TODO wrap custom typed error
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	return manifest, nil
}

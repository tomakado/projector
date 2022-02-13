package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "projector",
	Short: "Projector is a tool for generating project from templates",
	Long: `A flexible, language and framework agnostic tool that allows you to generate projects from templates. 
Projector has some builtin templates, but you can use your custom templates or import third-party templates
from GitHub.`,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// Execute runs passed command and handles errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

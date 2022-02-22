package cmd

import (
	"embed"
	"os"

	"github.com/spf13/cobra"
	"github.com/tomakado/projector/internal/pkg/verbose"
)

const embedRoot = "resources/templates/"

var (
	rootCmd = &cobra.Command{
		Use:   "projector",
		Short: "Projector is a tool for generating project from templates",
		Long: `A flexible, language and framework agnostic tool that allows you to generate projects from templates. 
Projector has some builtin templates, but you can use your custom templates or import third-party templates
from GitHub.`,
		SilenceUsage: true,
		PersistentPreRun: func(*cobra.Command, []string) {
			verbose.SetVerboseOn(isVerboseOn)
			verbose.Println("verbose mode is turned on")
		},
	}
	isVerboseOn bool
)

//go:embed resources/*
var resources embed.FS

type provider interface {
	Get(filename string) ([]byte, error)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&isVerboseOn, "verbose", "v", false, "turn verbose mode on")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs passed command and handles errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

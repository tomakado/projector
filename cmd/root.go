package cmd

import (
	"fmt"
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
	PersistentPostRun: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	// rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
	// 	cmd.Println("ohhh fuck")
	// 	// cmd.Printf("Error occurred during execution: %v\n", err)
	// 	return err
	// })

	rootCmd.AddCommand(createCmd)
}

// Execute runs passed command and handles errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("oh fuck")
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/tomakado/projector/internal/build"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display projector version",
	Run:   runVersion,
}

func runVersion(*cobra.Command, []string) {
	fmt.Println("projector")
	fmt.Println("===>{...}")
	fmt.Printf(
		"\n%s built at %s (commit %s on branch %q)\n",
		build.Version(),
		build.Time(),
		build.Commit(),
		build.Branch(),
	)
	fmt.Printf("%s@%s\n", runtime.GOOS, runtime.GOARCH)
}

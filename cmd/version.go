package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ghostchar version %s (built %s, commit %s)\n", version, date, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

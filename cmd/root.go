package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghostchar [flags] [path...]",
	Short: "Detect invisible Unicode, PUA, and bidi characters in source files",
	Long: `ghostchar scans source files for invisible Unicode characters, Private Use Area
codepoints, and bidirectional control characters. It reports exact locations
(file:line:column) and supports text, JSON, and SARIF output formats.

When called without a subcommand, behaves as "ghostchar scan".`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	// If no subcommand is matched, default to scan.
	// We intercept by inserting "scan" when the first arg is not a known command.
	cmd, _, _ := rootCmd.Find(os.Args[1:])
	if cmd == rootCmd && len(os.Args) > 1 {
		// No subcommand matched — prepend "scan" so cobra routes to scanCmd
		args := append([]string{"scan"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

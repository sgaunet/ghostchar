package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/sgaunet/ghostchar/internal/charset"
	"github.com/sgaunet/ghostchar/internal/report"
	"github.com/sgaunet/ghostchar/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	flagFormat      string
	flagQuiet       bool
	flagExt         string
	flagExclude     string
	flagCategories  string
	flagMaxFileSize string
	flagNoColor     bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [flags] [path...]",
	Short: "Scan files for invisible Unicode, PUA, and bidi characters",
	Long: `Scan source files for invisible Unicode characters, Private Use Area codepoints,
and bidirectional control characters. Reports findings with exact file:line:column locations.`,
	RunE: runScan,
}

func init() {
	scanCmd.Flags().StringVar(&flagFormat, "format", "text", "output format: text, json, sarif")
	scanCmd.Flags().BoolVarP(&flagQuiet, "quiet", "q", false, "suppress all output, only set exit code")
	scanCmd.Flags().StringVar(&flagExt, "ext", "", "comma-separated file extensions to scan (overrides defaults)")
	scanCmd.Flags().StringVar(&flagExclude, "exclude", "", "comma-separated directory names to exclude (overrides defaults)")
	scanCmd.Flags().StringVar(&flagCategories, "categories", "all", "comma-separated categories: invisible,pua,bidi,all")
	scanCmd.Flags().StringVar(&flagMaxFileSize, "max-file-size", "1MB", "maximum file size (e.g., 512KB, 2MB)")
	scanCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "disable colored output")
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	// Validate format flag
	switch flagFormat {
	case "text", "json", "sarif":
	default:
		return fmt.Errorf("invalid format %q: must be text, json, or sarif", flagFormat)
	}

	// Parse extensions
	extensions := scanner.DefaultExtensions()
	if flagExt != "" {
		extensions = make(map[string]bool)
		for _, ext := range strings.Split(flagExt, ",") {
			ext = strings.TrimSpace(ext)
			if ext != "" {
				extensions[ext] = true
			}
		}
	}

	// Parse exclude dirs
	excludeDirs := scanner.DefaultExcludeDirs()
	if flagExclude != "" {
		excludeDirs = make(map[string]bool)
		for _, dir := range strings.Split(flagExclude, ",") {
			dir = strings.TrimSpace(dir)
			if dir != "" {
				excludeDirs[dir] = true
			}
		}
	}

	// Parse categories
	var categories []charset.Category
	if flagCategories != "all" {
		for _, cat := range strings.Split(flagCategories, ",") {
			cat = strings.TrimSpace(cat)
			switch charset.Category(cat) {
			case charset.Invisible, charset.PUA, charset.Bidi:
				categories = append(categories, charset.Category(cat))
			case "all":
				categories = nil
				break
			default:
				return fmt.Errorf("invalid category %q: must be invisible, pua, bidi, or all", cat)
			}
		}
	}

	// Parse max file size
	maxFileSize, err := parseSize(flagMaxFileSize)
	if err != nil {
		return fmt.Errorf("invalid max-file-size %q: %w", flagMaxFileSize, err)
	}

	if flagNoColor {
		color.NoColor = true
	}

	cfg := scanner.ScanConfig{
		Paths:       args,
		Extensions:  extensions,
		ExcludeDirs: excludeDirs,
		Categories:  categories,
		MaxFileSize: maxFileSize,
		Format:      flagFormat,
		Quiet:       flagQuiet,
		NoColor:     flagNoColor,
	}

	warn := func(format string, args ...any) {
		fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
	}

	result, err := scanner.ScanPaths(cfg, warn)
	if err != nil {
		return err
	}

	if !flagQuiet {
		switch flagFormat {
		case "json":
			if err := report.FormatJSON(os.Stdout, result); err != nil {
				return err
			}
		case "sarif":
			if err := report.FormatSARIF(os.Stdout, result); err != nil {
				return err
			}
		default:
			report.FormatText(os.Stdout, result)
		}
	}

	if result.TotalFindings > 0 {
		os.Exit(1)
	}
	return nil
}

// parseSize parses a human-readable size string like "1MB", "512KB" into bytes.
func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	var multiplier int64 = 1
	if strings.HasSuffix(s, "GB") {
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "GB")
	} else if strings.HasSuffix(s, "MB") {
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "MB")
	} else if strings.HasSuffix(s, "KB") {
		multiplier = 1024
		s = strings.TrimSuffix(s, "KB")
	} else if strings.HasSuffix(s, "B") {
		s = strings.TrimSuffix(s, "B")
	}

	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, fmt.Errorf("size must be positive")
	}
	return n * multiplier, nil
}

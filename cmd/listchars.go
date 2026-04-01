package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sgaunet/ghostchar/internal/charset"
	"github.com/spf13/cobra"
)

var flagListCategory string

var listCharsCmd = &cobra.Command{
	Use:   "list-chars",
	Short: "List all detectable Unicode characters",
	Long:  "Display all Unicode characters and ranges that ghostchar can detect, grouped by category.",
	RunE:  runListChars,
}

func init() {
	listCharsCmd.Flags().StringVar(&flagListCategory, "category", "all", "filter by category: invisible, pua, bidi, all")
	rootCmd.AddCommand(listCharsCmd)
}

func runListChars(cmd *cobra.Command, args []string) error {
	var categories []charset.Category
	if flagListCategory == "all" {
		categories = charset.AllCategories()
	} else {
		cat := charset.Category(flagListCategory)
		switch cat {
		case charset.Invisible, charset.PUA, charset.Bidi:
			categories = []charset.Category{cat}
		default:
			return fmt.Errorf("invalid category %q: must be invisible, pua, bidi, or all", flagListCategory)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, cat := range categories {
		fmt.Fprintf(w, "\n[%s]\n", cat)
		fmt.Fprintf(w, "Codepoint\tName\n")
		fmt.Fprintf(w, "---------\t----\n")
		defs := charset.CharDefsByCategory(cat)
		for _, d := range defs {
			if d.IsSingle() {
				fmt.Fprintf(w, "U+%04X\t%s\n", d.Low, d.Name)
			} else {
				fmt.Fprintf(w, "U+%04X..U+%04X\t%s\n", d.Low, d.High, d.Name)
			}
		}
	}
	w.Flush()
	return nil
}

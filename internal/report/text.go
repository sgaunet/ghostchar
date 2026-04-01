package report

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/sgaunet/ghostchar/internal/charset"
	"github.com/sgaunet/ghostchar/internal/scanner"
)

var (
	colorBidi      = color.New(color.FgRed, color.Bold)
	colorInvisible = color.New(color.FgYellow)
	colorPUA       = color.New(color.FgCyan)
)

func colorForCategory(cat charset.Category) *color.Color {
	switch cat {
	case charset.Bidi:
		return colorBidi
	case charset.Invisible:
		return colorInvisible
	case charset.PUA:
		return colorPUA
	default:
		return color.New(color.Reset)
	}
}

// FormatText writes the scan result in human-readable text format.
func FormatText(w io.Writer, result *scanner.ScanResult) {
	for _, f := range result.Findings {
		c := colorForCategory(f.Category)
		fmt.Fprintf(w, "%s:%d:%d\t%s\t%s\t%s\n",
			f.File, f.Line, f.Column,
			c.Sprint(f.CodepointHex()), f.Name, c.Sprintf("[%s]", f.Category))
	}
	if len(result.Findings) > 0 {
		fmt.Fprintln(w)
	}
	fmt.Fprintf(w, "%d findings in %d file(s) (%d files scanned)\n",
		result.TotalFindings, result.FilesWithFindings, result.FilesScanned)
}

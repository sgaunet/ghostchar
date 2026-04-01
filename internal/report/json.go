package report

import (
	"encoding/json"
	"io"

	"github.com/sgaunet/ghostchar/internal/scanner"
)

type jsonOutput struct {
	Summary  jsonSummary   `json:"summary"`
	Findings []jsonFinding `json:"findings"`
}

type jsonSummary struct {
	FilesScanned      int `json:"files_scanned"`
	FilesWithFindings int `json:"files_with_findings"`
	TotalFindings     int `json:"total_findings"`
}

type jsonFinding struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Codepoint string `json:"codepoint"`
	Name      string `json:"name"`
	Category  string `json:"category"`
}

// FormatJSON writes the scan result as JSON.
func FormatJSON(w io.Writer, result *scanner.ScanResult) error {
	out := jsonOutput{
		Summary: jsonSummary{
			FilesScanned:      result.FilesScanned,
			FilesWithFindings: result.FilesWithFindings,
			TotalFindings:     result.TotalFindings,
		},
		Findings: make([]jsonFinding, 0, len(result.Findings)),
	}

	for _, f := range result.Findings {
		out.Findings = append(out.Findings, jsonFinding{
			File:      f.File,
			Line:      f.Line,
			Column:    f.Column,
			Codepoint: f.CodepointHex(),
			Name:      f.Name,
			Category:  string(f.Category),
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

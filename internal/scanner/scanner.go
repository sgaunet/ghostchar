package scanner

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"unicode/utf8"

	"github.com/sgaunet/ghostchar/internal/charset"
)

// Finding represents a single detected suspicious character occurrence.
type Finding struct {
	File      string           `json:"file"`
	Line      int              `json:"line"`
	Column    int              `json:"column"`
	Codepoint rune             `json:"-"`
	Name      string           `json:"name"`
	Category  charset.Category `json:"category"`
}

// CodepointHex returns the codepoint as a U+XXXX string.
func (f Finding) CodepointHex() string {
	if f.Codepoint > 0xFFFF {
		return fmt.Sprintf("U+%05X", f.Codepoint)
	}
	return fmt.Sprintf("U+%04X", f.Codepoint)
}

// ScanResult holds the aggregate output of a scan operation.
type ScanResult struct {
	FilesScanned      int       `json:"files_scanned"`
	FilesWithFindings int       `json:"files_with_findings"`
	TotalFindings     int       `json:"total_findings"`
	Findings          []Finding `json:"findings"`
}

// ScanConfig holds runtime configuration for a scan operation.
type ScanConfig struct {
	Paths       []string
	Extensions  map[string]bool
	ExcludeDirs map[string]bool
	Categories  []charset.Category
	Format      string
	Quiet       bool
	NoColor     bool
	MaxFileSize int64
}

// WarnFunc is a callback for reporting warnings (e.g., invalid UTF-8).
type WarnFunc func(format string, args ...any)

// ScanFile scans a single file for suspicious Unicode characters.
func ScanFile(path string, relPath string, categories []charset.Category, warn WarnFunc) ([]Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	lines := bytes.Split(data, []byte("\n"))

	detectFunc := charset.Detect
	if len(categories) > 0 && len(categories) < len(charset.AllCategories()) {
		detectFunc = func(r rune) (charset.CharDef, bool) {
			return charset.DetectInCategory(r, categories)
		}
	}

	for lineIdx, line := range lines {
		// Strip trailing \r for Windows line endings
		line = bytes.TrimSuffix(line, []byte("\r"))
		col := 1 // 1-based byte offset
		remaining := line
		for len(remaining) > 0 {
			r, size := utf8.DecodeRune(remaining)
			if r == utf8.RuneError && size == 1 {
				if warn != nil {
					warn("%s:%d:%d: invalid UTF-8 byte 0x%02X", relPath, lineIdx+1, col, remaining[0])
				}
				remaining = remaining[size:]
				col += size
				continue
			}

			if cd, ok := detectFunc(r); ok {
				name := cd.Name
				if !cd.IsSingle() {
					name = fmt.Sprintf("%s (U+%04X)", cd.Name, r)
				}
				findings = append(findings, Finding{
					File:      relPath,
					Line:      lineIdx + 1,
					Column:    col,
					Codepoint: r,
					Name:      name,
					Category:  cd.Category,
				})
			}
			remaining = remaining[size:]
			col += size
		}
	}
	return findings, nil
}

// ScanPaths scans the given paths (or cwd if empty) and returns aggregated results.
func ScanPaths(cfg ScanConfig, warn WarnFunc) (*ScanResult, error) {
	paths := cfg.Paths
	if len(paths) == 0 {
		paths = []string{"."}
	}

	files, err := WalkPaths(paths, cfg.Extensions, cfg.ExcludeDirs, cfg.MaxFileSize)
	if err != nil {
		return nil, err
	}

	result := &ScanResult{}
	for _, f := range files {
		result.FilesScanned++
		findings, err := ScanFile(f.AbsPath, f.RelPath, cfg.Categories, warn)
		if err != nil {
			if warn != nil {
				warn("error scanning %s: %v", f.RelPath, err)
			}
			continue
		}
		if len(findings) > 0 {
			result.FilesWithFindings++
			result.TotalFindings += len(findings)
			result.Findings = append(result.Findings, findings...)
		}
	}

	// Sort findings by file, then line, then column
	sort.Slice(result.Findings, func(i, j int) bool {
		fi, fj := result.Findings[i], result.Findings[j]
		if fi.File != fj.File {
			return fi.File < fj.File
		}
		if fi.Line != fj.Line {
			return fi.Line < fj.Line
		}
		return fi.Column < fj.Column
	})

	return result, nil
}

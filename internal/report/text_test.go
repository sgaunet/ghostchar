package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sgaunet/ghostchar/internal/charset"
	"github.com/sgaunet/ghostchar/internal/scanner"
)

func TestFormatText_WithFindings(t *testing.T) {
	result := &scanner.ScanResult{
		FilesScanned:      5,
		FilesWithFindings: 1,
		TotalFindings:     1,
		Findings: []scanner.Finding{
			{
				File: "test.go", Line: 3, Column: 5,
				Codepoint: 0x200B, Name: "ZERO WIDTH SPACE", Category: charset.Invisible,
			},
		},
	}

	var buf bytes.Buffer
	FormatText(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "test.go:3:5") {
		t.Errorf("expected file:line:column, got:\n%s", out)
	}
	if !strings.Contains(out, "U+200B") {
		t.Errorf("expected U+200B in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[invisible]") {
		t.Errorf("expected [invisible] in output, got:\n%s", out)
	}
	if !strings.Contains(out, "1 findings in 1 file(s) (5 files scanned)") {
		t.Errorf("expected summary line, got:\n%s", out)
	}
}

func TestFormatText_Empty(t *testing.T) {
	result := &scanner.ScanResult{FilesScanned: 3}
	var buf bytes.Buffer
	FormatText(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "0 findings in 0 file(s) (3 files scanned)") {
		t.Errorf("expected empty summary, got:\n%s", out)
	}
}

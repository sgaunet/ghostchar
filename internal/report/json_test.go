package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sgaunet/ghostchar/internal/charset"
	"github.com/sgaunet/ghostchar/internal/scanner"
)

func TestFormatJSON_EmptyResult(t *testing.T) {
	result := &scanner.ScanResult{FilesScanned: 5}
	var buf bytes.Buffer
	if err := FormatJSON(&buf, result); err != nil {
		t.Fatal(err)
	}

	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Summary.FilesScanned != 5 {
		t.Errorf("expected files_scanned=5, got %d", out.Summary.FilesScanned)
	}
	if len(out.Findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(out.Findings))
	}
}

func TestFormatJSON_WithFindings(t *testing.T) {
	result := &scanner.ScanResult{
		FilesScanned:      10,
		FilesWithFindings: 1,
		TotalFindings:     1,
		Findings: []scanner.Finding{
			{
				File:      "test.go",
				Line:      3,
				Column:    5,
				Codepoint: 0x200B,
				Name:      "ZERO WIDTH SPACE",
				Category:  charset.Invisible,
			},
		},
	}
	var buf bytes.Buffer
	if err := FormatJSON(&buf, result); err != nil {
		t.Fatal(err)
	}

	var out jsonOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Summary.TotalFindings != 1 {
		t.Errorf("expected total_findings=1, got %d", out.Summary.TotalFindings)
	}
	if out.Findings[0].Codepoint != "U+200B" {
		t.Errorf("expected codepoint U+200B, got %s", out.Findings[0].Codepoint)
	}
	if out.Findings[0].Category != "invisible" {
		t.Errorf("expected category invisible, got %s", out.Findings[0].Category)
	}
}

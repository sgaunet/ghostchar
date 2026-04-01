package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sgaunet/ghostchar/internal/charset"
	"github.com/sgaunet/ghostchar/internal/scanner"
)

func TestFormatSARIF_Structure(t *testing.T) {
	result := &scanner.ScanResult{
		FilesScanned:      1,
		FilesWithFindings: 1,
		TotalFindings:     2,
		Findings: []scanner.Finding{
			{
				File: "test.go", Line: 1, Column: 5,
				Codepoint: 0x200B, Name: "ZERO WIDTH SPACE", Category: charset.Invisible,
			},
			{
				File: "test.go", Line: 2, Column: 10,
				Codepoint: 0x202E, Name: "RIGHT-TO-LEFT OVERRIDE", Category: charset.Bidi,
			},
		},
	}

	var buf bytes.Buffer
	if err := FormatSARIF(&buf, result); err != nil {
		t.Fatal(err)
	}

	var log sarifLog
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if log.Version != "2.1.0" {
		t.Errorf("expected version 2.1.0, got %s", log.Version)
	}
	if len(log.Runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(log.Runs))
	}

	run := log.Runs[0]
	if run.Tool.Driver.Name != "ghostchar" {
		t.Errorf("expected tool name ghostchar, got %s", run.Tool.Driver.Name)
	}
	if len(run.Tool.Driver.Rules) != 3 {
		t.Errorf("expected 3 rules, got %d", len(run.Tool.Driver.Rules))
	}
	if len(run.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(run.Results))
	}

	// Check invisible finding
	r0 := run.Results[0]
	if r0.RuleID != "ghostchar/invisible" {
		t.Errorf("expected ruleId ghostchar/invisible, got %s", r0.RuleID)
	}
	if r0.Level != "warning" {
		t.Errorf("expected level warning for invisible, got %s", r0.Level)
	}

	// Check bidi finding
	r1 := run.Results[1]
	if r1.RuleID != "ghostchar/bidi" {
		t.Errorf("expected ruleId ghostchar/bidi, got %s", r1.RuleID)
	}
	if r1.Level != "error" {
		t.Errorf("expected level error for bidi, got %s", r1.Level)
	}
	if r1.Locations[0].PhysicalLocation.Region.StartLine != 2 {
		t.Errorf("expected line 2, got %d", r1.Locations[0].PhysicalLocation.Region.StartLine)
	}
}

func TestFormatSARIF_EmptyResult(t *testing.T) {
	result := &scanner.ScanResult{FilesScanned: 3}
	var buf bytes.Buffer
	if err := FormatSARIF(&buf, result); err != nil {
		t.Fatal(err)
	}

	var log sarifLog
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(log.Runs[0].Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(log.Runs[0].Results))
	}
}

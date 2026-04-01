package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sgaunet/ghostchar/internal/charset"
)

func TestScanFile_Invisible(t *testing.T) {
	// Create a temp file with a ZWSP
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	// "hello\u200bworld" on line 1
	os.WriteFile(path, []byte("hello\xe2\x80\x8bworld\n"), 0o644)

	findings, err := ScanFile(path, "test.go", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	f := findings[0]
	if f.Codepoint != 0x200B {
		t.Errorf("expected U+200B, got U+%04X", f.Codepoint)
	}
	if f.Line != 1 {
		t.Errorf("expected line 1, got %d", f.Line)
	}
	if f.Column != 6 {
		t.Errorf("expected column 6, got %d", f.Column)
	}
	if f.Category != charset.Invisible {
		t.Errorf("expected invisible, got %s", f.Category)
	}
}

func TestScanFile_MultipleFindingsOneLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	// Two ZWSP characters: positions 1 and 5 (1-based byte offset)
	os.WriteFile(path, []byte("\xe2\x80\x8babc\xe2\x80\x8b\n"), 0o644)

	findings, err := ScanFile(path, "test.go", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	if findings[0].Column != 1 {
		t.Errorf("first finding: expected column 1, got %d", findings[0].Column)
	}
	if findings[1].Column != 7 {
		t.Errorf("second finding: expected column 7, got %d", findings[1].Column)
	}
}

func TestScanFile_CategoryFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	// ZWSP (invisible) + RLO (bidi)
	os.WriteFile(path, []byte("\xe2\x80\x8b\xe2\x80\xae\n"), 0o644)

	// Filter for bidi only
	findings, err := ScanFile(path, "test.go", []charset.Category{charset.Bidi}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding (bidi only), got %d", len(findings))
	}
	if findings[0].Category != charset.Bidi {
		t.Errorf("expected bidi, got %s", findings[0].Category)
	}
}

func TestScanFile_Clean(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "clean.go")
	os.WriteFile(path, []byte("package main\n\nfunc main() {}\n"), 0o644)

	findings, err := ScanFile(path, "clean.go", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestScanFile_InvalidUTF8Warning(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	os.WriteFile(path, []byte("hello\x80world\n"), 0o644)

	var warnings []string
	warn := func(format string, args ...any) {
		warnings = append(warnings, "warned")
	}

	_, err := ScanFile(path, "test.go", nil, warn)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) == 0 {
		t.Error("expected warning for invalid UTF-8")
	}
}

func TestCodepointHex(t *testing.T) {
	tests := []struct {
		cp   rune
		want string
	}{
		{0x200B, "U+200B"},
		{0x00AD, "U+00AD"},
		{0xE000, "U+E000"},
		{0xF0000, "U+F0000"},
	}
	for _, tt := range tests {
		f := Finding{Codepoint: tt.cp}
		got := f.CodepointHex()
		if got != tt.want {
			t.Errorf("CodepointHex(%#U) = %s, want %s", tt.cp, got, tt.want)
		}
	}
}

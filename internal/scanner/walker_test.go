package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkPaths_FiltersExtensions(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("package main\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "test.xyz"), []byte("ignored\n"), 0o644)

	files, err := WalkPaths([]string{dir}, map[string]bool{"go": true}, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if filepath.Ext(files[0].AbsPath) != ".go" {
		t.Errorf("expected .go file, got %s", files[0].AbsPath)
	}
}

func TestWalkPaths_ExcludesDirs(t *testing.T) {
	dir := t.TempDir()
	vendorDir := filepath.Join(dir, "vendor")
	os.MkdirAll(vendorDir, 0o755)
	os.WriteFile(filepath.Join(vendorDir, "lib.go"), []byte("package vendor\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644)

	files, err := WalkPaths([]string{dir}, nil, map[string]bool{"vendor": true}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file (vendor excluded), got %d", len(files))
	}
}

func TestWalkPaths_SkipsBinary(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "binary.go"), append([]byte("package main\n\x00binary data"), make([]byte, 100)...), 0o644)
	os.WriteFile(filepath.Join(dir, "text.go"), []byte("package main\n"), 0o644)

	files, err := WalkPaths([]string{dir}, map[string]bool{"go": true}, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file (binary skipped), got %d", len(files))
	}
}

func TestWalkPaths_MaxFileSize(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "small.go"), []byte("package main\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "big.go"), make([]byte, 2000), 0o644)

	files, err := WalkPaths([]string{dir}, map[string]bool{"go": true}, nil, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file (big skipped), got %d", len(files))
	}
}

func TestWalkPaths_SingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	os.WriteFile(path, []byte("package main\n"), 0o644)

	files, err := WalkPaths([]string{path}, nil, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
}

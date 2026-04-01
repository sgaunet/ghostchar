package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileEntry represents a file to be scanned.
type FileEntry struct {
	AbsPath string
	RelPath string
}

// DefaultExtensions returns the default set of file extensions to scan.
func DefaultExtensions() map[string]bool {
	exts := []string{
		"go", "py", "js", "ts", "jsx", "tsx", "java", "c", "cpp", "h", "hpp",
		"cs", "rb", "php", "rs", "kt", "swift", "sh", "bash", "zsh",
		"yaml", "yml", "toml", "json", "xml", "html", "htm", "css", "scss",
		"sql", "tf", "md", "txt",
	}
	m := make(map[string]bool, len(exts))
	for _, e := range exts {
		m[e] = true
	}
	return m
}

// DefaultExcludeDirs returns the default set of directory names to exclude.
func DefaultExcludeDirs() map[string]bool {
	return map[string]bool{
		".git":         true,
		"vendor":       true,
		"node_modules": true,
	}
}

// WalkPaths walks the given paths and collects files matching the filter criteria.
func WalkPaths(paths []string, extensions map[string]bool, excludeDirs map[string]bool, maxFileSize int64) ([]FileEntry, error) {
	if extensions == nil {
		extensions = DefaultExtensions()
	}
	if excludeDirs == nil {
		excludeDirs = DefaultExcludeDirs()
	}
	if maxFileSize <= 0 {
		maxFileSize = 1048576 // 1 MB default
	}

	var files []FileEntry

	for _, root := range paths {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			return nil, err
		}

		info, err := os.Stat(absRoot)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			// Single file
			rel, _ := filepath.Rel(".", root)
			if rel == "" {
				rel = root
			}
			files = append(files, FileEntry{AbsPath: absRoot, RelPath: rel})
			continue
		}

		err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // skip entries with errors
			}

			// Skip symlinks
			if d.Type()&os.ModeSymlink != 0 {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if d.IsDir() {
				name := d.Name()
				if excludeDirs[name] {
					return filepath.SkipDir
				}
				return nil
			}

			// Check extension
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			if !extensions[ext] {
				return nil
			}

			// Check file size
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if info.Size() > maxFileSize {
				return nil
			}

			// Check for binary content (null byte in first 512 bytes)
			if isBinary(path) {
				return nil
			}

			rel, err := filepath.Rel(absRoot, path)
			if err != nil {
				rel = path
			}
			// Prefix with the original root path for relative display
			if root != "." {
				rel = filepath.Join(root, rel)
			}

			files = append(files, FileEntry{AbsPath: path, RelPath: rel})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

// isBinary checks if a file appears to be binary by looking for null bytes in the first 512 bytes.
func isBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if n == 0 {
		return false
	}
	for _, b := range buf[:n] {
		if b == 0 {
			return true
		}
	}
	return false
}

package db

import (
	"embed"
	"path/filepath"
)

//go:embed migrations
var migrationFS embed.FS

func MigrationsFS() embed.FS {
	return migrationFS
}

func MigrationsBaseDir() string {
	return filepath.Join("migrations")
}

func EntryNames(fs embed.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}

	return names, nil
}

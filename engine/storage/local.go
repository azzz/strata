package storage

import (
	"context"
	"fmt"
	"path/filepath"
)

// LocalStorage is a storage implementation that reads data from the local filesystem.
type LocalStorage struct {
	RootDir string
}

// NewLocalStorage creates a new LocalStorage instance with the specified root directory.
func (s *LocalStorage) NewScanner(ctx context.Context, req ScanRequest) (Scanner, error) {
	switch req.Format {
	case FormatParquet:
		// Implement Parquet scanner creation here
		return nil, fmt.Errorf("parquet format not implemented yet")
	case FormatJSONL:
		path := filepath.Join(s.RootDir, string(req.URI))
		scanner := NewJSONLScanner(path, req.Schema)

		return scanner, nil
	}

	return nil, nil
}

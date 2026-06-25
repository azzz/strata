package storage

import (
	"context"

	"github.com/azzz/strata/engine/types"
)

type URI string
type Value any
type Format string

const (
	FormatParquet Format = "parquet"
	FormatJSONL   Format = "jsonl"
)

type Scanner interface {
	Next(ctx context.Context) bool
	Err() error
	Row() types.Row
	Close() error
	Open() error
}

// ScanRequest represents a request to scan data from a specific URI with a given format and schema.
type ScanRequest struct {
	URI    URI
	Format Format
	Schema types.Schema
}

type Storage interface {
	NewScanner(ctx context.Context, req ScanRequest) (Scanner, error)
}

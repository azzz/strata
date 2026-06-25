package storage

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/azzz/strata/engine/storage/jsonl"
	"github.com/azzz/strata/engine/types"
)

// JSONLScanner is a scanner for reading JSONL files.
// It reads the file line by line and decodes each line into a Row according to the provided schema.
type JSONLScanner struct {
	path   string
	schema types.Schema

	file    *os.File
	decoder *jsonl.Decoder
	scanner *bufio.Scanner

	row      types.Row
	rowIndex int
	err      error
}

func NewJSONLScanner(path string, schema types.Schema) *JSONLScanner {
	return &JSONLScanner{
		path:   path,
		schema: schema,
	}
}

// Open opens the JSONL file for reading and initializes the scanner and decoder.
func (s *JSONLScanner) Open() error {
	file, err := os.Open(string(s.path))
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	s.file = file
	s.scanner = bufio.NewScanner(file)
	s.rowIndex = 0

	s.decoder = jsonl.NewDecoder(s.schema)

	return nil
}

// Next reads the next line from the JSONL file, decodes it into a Row, and returns true if successful.
func (s *JSONLScanner) Next(ctx context.Context) bool {
	if s.file == nil || s.scanner == nil {
		s.err = fmt.Errorf("scanner is not opened")
		return false
	}

	select {
	case <-ctx.Done():
		s.err = ctx.Err()
		return false
	default:
	}

	if !s.scanner.Scan() {
		if err := s.scanner.Err(); err != nil {
			s.err = fmt.Errorf("failed to read next line: %w", err)
		}

		return false
	}

	line := s.scanner.Bytes()

	row, err := s.decoder.Decode(line)
	if err != nil {
		s.err = fmt.Errorf("failed to decode row: %w", err)
		return false
	}

	s.row = row
	s.rowIndex++

	return true
}

// Err returns the last error encountered by the scanner.
func (s *JSONLScanner) Err() error {
	return s.err
}

// Row returns the current row read by the scanner.
func (s *JSONLScanner) Row() types.Row {
	return s.row
}

// Close closes the JSONL file and releases any resources associated with the scanner.
func (s *JSONLScanner) Close() error {
	if s.file == nil {
		return nil
	}

	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close file handler: %w", err)
	}

	return nil
}

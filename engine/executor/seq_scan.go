package executor

import (
	"context"
	"fmt"

	"github.com/azzz/strata/engine/storage"
	"github.com/azzz/strata/engine/types"
)

type SeqScan struct {
	storage  storage.Storage
	requests []storage.ScanRequest

	current storage.Scanner
	idx     int

	row types.Row
	err error
}

func NewSeqScan(storage storage.Storage, requests []storage.ScanRequest) *SeqScan {
	return &SeqScan{
		storage:  storage,
		requests: requests,
	}
}

func (s *SeqScan) Next(ctx context.Context) bool {
	for s.idx < len(s.requests) {
		if s.current == nil {
			scanner, err := s.storage.NewScanner(ctx, s.requests[s.idx])
			if err != nil {
				s.err = fmt.Errorf("failed to open scanner (%d): %w", s.idx, err)
				return false
			}

			if err := scanner.Open(); err != nil {
				s.err = fmt.Errorf("failed to open scanner (%d): %w", s.idx, err)
				return false
			}

			s.current = scanner
		}

		// everything is good
		if s.current.Next(ctx) {
			s.row = s.current.Row()
			return true
		}

		// handle current scanner error
		if err := s.current.Err(); err != nil {
			s.err = fmt.Errorf("scanner failed (%d): %w", s.idx, err)
			return false
		}

		// no error? Scanner finished correctly => close current scanner
		if err := s.current.Close(); err != nil {
			s.err = fmt.Errorf("failed to close current scanner (%d): %w", s.idx, err)
			return false
		}

		// update cursor
		s.idx++
		s.current = nil
	}

	return false
}

func (r *SeqScan) Row() types.Row { return r.row }
func (r *SeqScan) Err() error     { return r.err }

func (r *SeqScan) Close() error {
	return r.current.Close()
}

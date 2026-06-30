package executor

import (
	"context"

	"github.com/azzz/strata/engine/plan"
	"github.com/azzz/strata/engine/storage"
	"github.com/azzz/strata/engine/types"
)

// Executor builds execution pipeline using Plan.
// Basic execution pipeline looks like:
// Limit => Offset => Filter => SeqScan => StorageScanner
type Executor struct {
	Plan    plan.PhysicalPlan
	Schema  types.Schema
	Storage storage.Storage
}

func (e Executor) Exec(ctx context.Context) Operator {
	requests := make([]storage.ScanRequest, 0, len(e.Plan.Scans))

	for _, scan := range e.Plan.Scans {
		req := storage.ScanRequest{
			URI:     storage.URI(scan.URI),
			Format:  scan.Format,
			Schema:  e.Schema,
			Columns: e.Plan.RequestedColumns,
		}

		requests = append(requests, req)
	}

	var root Operator = &SeqScan{
		requests: requests,
		storage:  e.Storage,
	}

	if e.Plan.Filter != nil {
		root = NewFilter(root, e.Plan.Filter)
	}

	if e.Plan.Offset > 0 {
		root = NewOffset(root, e.Plan.Offset)
	}

	if e.Plan.Limit > 0 {
		root = NewLimit(root, e.Plan.Limit)
	}

	return root
}

package executor

import (
	"context"

	"github.com/azzz/strata/engine/types"
)

type Operator interface {
	Next(ctx context.Context) bool
	Row() types.Row
	Err() error
	Close() error
}

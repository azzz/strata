package executor

import (
	"context"

	"github.com/azzz/strata/engine/types"
)

type Limit struct {
	input Operator

	limit   int
	emitted int

	row types.Row
	err error
}

func NewLimit(input Operator, limit int) *Limit {
	return &Limit{
		input: input,
		limit: limit,
	}
}
func (l *Limit) Row() types.Row { return l.row }
func (l *Limit) Err() error     { return l.err }
func (l *Limit) Close() error   { return l.input.Close() }

func (l *Limit) Next(ctx context.Context) bool {
	if l.limit > 0 && l.emitted >= l.limit {
		return false
	}

	if !l.input.Next(ctx) {
		l.err = l.input.Err()
		return false
	}

	l.row = l.input.Row()
	l.emitted++

	return true
}

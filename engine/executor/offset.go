package executor

import (
	"context"

	"github.com/azzz/strata/engine/types"
)

type Offset struct {
	input Operator

	offset  int
	skipped int

	row types.Row
	err error
}

func NewOffset(input Operator, offset int) *Offset {
	return &Offset{
		input:  input,
		offset: offset,
	}
}

func (o *Offset) Row() types.Row { return o.row }
func (o *Offset) Err() error     { return o.err }
func (o *Offset) Close() error   { return o.input.Close() }

func (o *Offset) Next(ctx context.Context) bool {
	for o.skipped < o.offset {
		if !o.input.Next(ctx) {
			o.err = o.input.Err()
			return false
		}

		o.skipped++
	}

	if !o.input.Next(ctx) {
		o.err = o.input.Err()
		return false
	}

	o.row = o.input.Row()

	return true
}

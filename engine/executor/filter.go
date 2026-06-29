package executor

import (
	"context"
	"fmt"

	"github.com/azzz/strata/engine/expr"
	"github.com/azzz/strata/engine/types"
)

type Filter struct {
	input  Operator
	filter expr.Filter

	row types.Row
	err error
}

func NewFilter(input Operator, filter expr.Filter) *Filter {
	return &Filter{
		input:  input,
		filter: filter,
	}
}

func (f *Filter) Row() types.Row { return f.row }
func (f *Filter) Err() error     { return f.err }
func (f *Filter) Close() error   { return f.input.Close() }

func (f *Filter) Next(ctx context.Context) bool {
	for f.input.Next(ctx) {
		row := f.input.Row()

		ok, err := f.filter.Match(row)
		if err != nil {
			f.err = fmt.Errorf("failed to apply filter: %w", err)
			return false
		}

		if !ok {
			continue
		}

		f.row = row
		return true
	}

	if err := f.input.Err(); err != nil {
		f.err = f.input.Err()
		return false
	}

	return false
}

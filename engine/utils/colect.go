package utils

import (
	"context"

	"github.com/azzz/strata/engine/types"
)

type RowScanner interface {
	Next(context.Context) bool
	Row() types.Row
	Err() error
}

func Collect(ctx context.Context, sc RowScanner) ([]types.Row, error) {
	result := make([]types.Row, 0, 100)
	for sc.Next(ctx) {
		result = append(result, sc.Row())
	}

	if err := sc.Err(); err != nil {
		return result, err
	}

	return result, nil
}

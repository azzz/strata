package expr

import "github.com/azzz/strata/engine/types"

type Eq struct {
	Col   int
	Value types.Value
}

func (e *Eq) Match(row types.Row) (bool, error) {
	result, err := compare(row, e.Col, e.Value)
	if err != nil {
		return false, err
	}

	return result == types.Equal, nil
}

type GreaterThan struct {
	Col   int
	Value types.Value
}

func (g *GreaterThan) Match(row types.Row) (bool, error) {
	result, err := compare(row, g.Col, g.Value)
	if err != nil {
		return false, err
	}

	return result == types.Greater, nil
}

type GreaterThanOrEqual struct {
	Col   int
	Value types.Value
}

func (g *GreaterThanOrEqual) Match(row types.Row) (bool, error) {
	result, err := compare(row, g.Col, g.Value)
	if err != nil {
		return false, err
	}

	return result >= types.Equal, nil
}

type LessThan struct {
	Col   int
	Value types.Value
}

func (l *LessThan) Match(row types.Row) (bool, error) {
	result, err := compare(row, l.Col, l.Value)
	if err != nil {
		return false, err
	}

	return result == types.Less, nil
}

type LessThanOrEqual struct {
	Col   int
	Value types.Value
}

func (l *LessThanOrEqual) Match(row types.Row) (bool, error) {
	result, err := compare(row, l.Col, l.Value)
	if err != nil {
		return false, err
	}

	return result <= types.Equal, nil
}

func compare(row types.Row, col int, value types.Value) (int, error) {
	v, ok := row.Get(col)
	if !ok {
		return 0, &MissingColumnError{Col: col}
	}

	return v.Compare(value)
}

package expr

import "github.com/azzz/strata/engine/types"

type And struct {
	Filters []Filter
}

func (a *And) Match(row types.Row) (bool, error) {
	for _, filter := range a.Filters {
		ok, err := filter.Match(row)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}

	return true, nil
}

type Or struct {
	Filters []Filter
}

func (o *Or) Match(row types.Row) (bool, error) {
	for _, filter := range o.Filters {
		ok, err := filter.Match(row)
		if err != nil {
			return false, err
		}

		if ok {
			return true, nil
		}
	}

	return false, nil
}

type Not struct {
	Filter Filter
}

func (n *Not) Match(row types.Row) (bool, error) {
	ok, err := n.Filter.Match(row)
	if err != nil {
		return false, err
	}

	return !ok, nil
}

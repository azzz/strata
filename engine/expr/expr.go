package expr

import "github.com/azzz/strata/engine/types"

type Filter interface {
	Match(types.Row) (bool, error)
}

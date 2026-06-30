package expr

import (
	"fmt"

	"github.com/azzz/strata/engine/types"
)

type MissingColumnError struct {
	Col types.ColumnIndex
}

func (e *MissingColumnError) Error() string {
	return "missing column: " + fmt.Sprint(e.Col)
}

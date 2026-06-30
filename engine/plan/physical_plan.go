package plan

import (
	"github.com/azzz/strata/engine/expr"
	"github.com/azzz/strata/engine/types"
)

type PhysicalPlan struct {
	Scans            []Scan
	Filter           expr.Filter
	Limit            int
	Offset           int
	RequestedColumns []types.ColumnIndex
}

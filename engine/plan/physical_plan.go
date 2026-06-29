package plan

import "github.com/azzz/strata/engine/expr"

type PhysicalPlan struct {
	Scans  []Scan
	Filter expr.Filter
	Limit  int
	Offset int
}

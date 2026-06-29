package plan

import "github.com/azzz/strata/engine/storage"

type Scan struct {
	URI     string
	Format  storage.Format
	Columns []int
}

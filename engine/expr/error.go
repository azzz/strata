package expr

import "fmt"

type MissingColumnError struct {
	Col int
}

func (e *MissingColumnError) Error() string {
	return "missing column: " + fmt.Sprint(e.Col)
}

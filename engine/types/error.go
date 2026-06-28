package types

import "fmt"

// TypeError represents an error that occurs when a value does not match the expected type.
type TypeError struct {
	Field    string // The name of the field that caused the error
	Expected Kind   // The expected type of the field
	Actual   Kind   // The actual type of the field
}

func (e *TypeError) Error() string {
	if e.Actual == "" {
		return "type error: field " + e.Field + " expected " + string(e.Expected) + " but got unknown type"
	}

	if e.Expected == "" && e.Actual == "" {
		return "type error: field " + e.Field + " is unknown type"
	}

	return "type error: field " + e.Field + " expected " + string(e.Expected) + " but got " + string(e.Actual)
}

type CompareError struct {
	Left  Kind
	Right Kind
}

func (e *CompareError) Error() string {
	return fmt.Sprintf("cannot compare %s with %s", e.Left, e.Right)
}

package types

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

	return "type error: field " + e.Field + " expected " + string(e.Expected) + " but got " + string(e.Actual)
}

// TypeNotImplementedError represents an error that occurs when a value's type is not implemented in the system.
type TypeNotImplementedError struct {
	Field string // The name of the field that caused the error
	Kind  Kind   // The kind of the field that is not implemented
}

func (e *TypeNotImplementedError) Error() string {
	return "type not implemented error: field " + e.Field + " of kind " + string(e.Kind) + " is not implemented"
}

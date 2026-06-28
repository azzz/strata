package types

const (
	// Less indicates that the left value is less than the right value.
	// Equal indicates that the left value is equal to the right value.
	// Greater indicates that the left value is greater than the right value.
	Less    = -1
	Equal   = 0
	Greater = 1
)

// Compare two values of the same type. Returns -1 if v < other, 0 if v == other, and 1 if v > other. If the types are different, returns an error.
func (v Value) Compare(other Value) (int, error) {
	if v.Kind != other.Kind {
		return 0, &CompareError{Left: v.Kind, Right: other.Kind}
	}

	switch v.Kind {
	case KindString:
		if v.S < other.S {
			return Less, nil
		} else if v.S > other.S {
			return Greater, nil
		}
		return Equal, nil
	case KindInt64:
		if v.I64 < other.I64 {
			return Less, nil
		} else if v.I64 > other.I64 {
			return Greater, nil
		}
		return Equal, nil
	case KindUInt64:
		if v.U64 < other.U64 {
			return Less, nil
		} else if v.U64 > other.U64 {
			return Greater, nil
		}
		return Equal, nil
	case KindFloat64:
		if v.F64 < other.F64 {
			return Less, nil
		} else if v.F64 > other.F64 {
			return Greater, nil
		}
		return Equal, nil
	case KindBool:
		if !v.B && other.B {
			return Less, nil
		} else if v.B && !other.B {
			return Greater, nil
		}
		return Equal, nil
	default:
		return 0, &CompareError{Left: v.Kind, Right: other.Kind}
	}
}

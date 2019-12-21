package gadget

import (
	"fmt"
	"go/parser"
	"reflect"
)

// ParseType takes a string containing a Go type definition,
// and returns the Type.
func ParseType(s string) (Type, error) {
	e, err := parser.ParseExpr(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %w", err)
	}
	t, err := convertTypeSpec(e)
	if err != nil {
		return nil, fmt.Errorf("failed to convert type: %w", err)
	}
	return t, nil
}

// TypeIs checks if the given type is equal to the type defined in the given string, as parsed by ParseType.
func TypeIs(t Type, s string) bool {
	t2, err := ParseType(s)
	if err != nil {
		return false
	}
	return SameType(t, t2)
}

// SameType will return true if the given types are the same.
func SameType(t, t2 Type) bool {
	return reflect.DeepEqual(t, t2)
}

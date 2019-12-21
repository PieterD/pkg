package gadget

import (
	"go/parser"
	"reflect"

	"github.com/pkg/errors"
)

// ParseType takes a string containing a Go type definition,
// and returns the Type.
func ParseType(s string) (Type, error) {
	e, err := parser.ParseExpr(s)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse expression")
	}
	t, err := convertTypeSpec(e)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert type")
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

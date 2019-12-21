package gadget

import "testing"

func TestTypeIs(t *testing.T) {
	for _, test := range []struct {
		str string
		typ Type
	}{
		{"string", String},
		{"bool", Bool},
		{"byte", Byte},
		{"rune", Rune},
		{"uintptr", Uintptr},
		{"int", Int},
		{"int8", Int8},
		{"int16", Int16},
		{"int32", Int32},
		{"int64", Int64},
		{"uint", Uint},
		{"uint8", Uint8},
		{"uint16", Uint16},
		{"uint32", Uint32},
		{"uint64", Uint64},
		{"float32", Float32},
		{"float64", Float64},
		{"complex64", Complex64},
		{"complex128", Complex128},
		{"struct{}", EmptyStruct},
		{"[]byte", Bytes},
		{"error", Error},
		{"io.Reader", Selector{Left: Ident("io"), Right: Ident("Reader")}},
	} {
		t.Run(test.str, func(t *testing.T) {
			if !TypeIs(test.typ, test.str) {
				t.Errorf("expected type string %s to be equal to type %#v", test.str, test.typ)
			}
			if test.typ.String() != test.str {
				t.Errorf("expected type %#v String (%s) to be equal to type string %s", test.typ, test.typ.String(), test.str)
			}
		})
	}
}

package gadget

import (
	"reflect"
	"testing"
)

func TestParseType(t *testing.T) {
	for _, test := range []struct {
		s string
		t Type
	}{
		{
			s: "SomethingMysterious",
			t: Ident("SomethingMysterious"),
		},
		{
			s: "os.Something",
			t: Selector{Left: "os", Right: "Something"},
		},
		{
			s: "[]byte",
			t: Bytes,
		},
		{
			s: "struct{}",
			t: EmptyStruct,
		},
		{
			s: "[16]bool",
			t: Array{Size: 16, Elem: Bool},
		},
		{
			s: "chan uint16",
			t: Chan{Dir: BOTH, Elem: Uint16},
		},
		{
			s: "<-chan bool",
			t: Chan{Dir: RECV, Elem: Bool},
		},
		{
			s: "chan<- [2]int",
			t: Chan{Dir: SEND, Elem: Array{Size: 2, Elem: Int}},
		},
		{
			s: "map[string][]*int",
			t: Map{Key: Ident("string"), Value: Slice{Elem: Pointer{Elem: Ident("int")}}},
		},
		{
			s: "struct {Embed; Name string; Age int `tag`}",
			t: Struct{Fields: []StructField{
				{Type: Ident("Embed")},
				{Name: "Name", Type: String},
				{Name: "Age", Type: Int, Tag: "tag"},
			}},
		},
		{
			s: "func()",
			t: Func{},
		},
		{
			s: "func([]string) []int",
			t: Func{
				Params: []FuncParam{
					{Type: Slice{Elem: String}},
				},
				Results: []FuncResult{
					{Type: Slice{Elem: Int}},
				},
			},
		},
		{
			s: "func(name, address string) (age, salary int)",
			t: Func{
				Params: []FuncParam{
					{Name: "name", Type: String},
					{Name: "address", Type: String},
				},
				Results: []FuncResult{
					{Name: "age", Type: Int},
					{Name: "salary", Type: Int},
				},
			},
		},
		{
			s: "interface{Get() int; Set(int)}",
			t: Interface{Methods: []InterfaceMethod{
				{Name: "Get", Type: Func{Results: []FuncResult{{Type: Int}}}},
				{Name: "Set", Type: Func{Params: []FuncParam{{Type: Int}}}},
			}},
		},
	} {
		t.Run(test.s, func(t *testing.T) {
			typ, err := ParseType(test.s)
			if err != nil {
				t.Errorf("Failed to parse: %v", err)
			}
			if !reflect.DeepEqual(typ, test.t) {
				t.Logf("want: %#v", test.t)
				t.Logf(" got: %#v", typ)
				t.Errorf("Types not equal")
			}
		})
	}
}

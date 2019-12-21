package gadget

import (
	"reflect"
	"testing"
)

func TestNewFile(t *testing.T) {
	path := "example/type.go"
	f, err := NewFile(path, nil)
	if err != nil {
		t.Fatalf("failed to create new file: %v", err)
	}

	expectedImports := []ImportDecl{
		{
			Position: Position{Path: path, Line: 4},
			Name:     "",
			Path:     "fmt",
		},
		{
			Position: Position{Path: path, Line: 5},
			Name:     "",
			Path:     "io",
		},
		{
			Position: Position{Path: path, Line: 6},
			Name:     ".",
			Path:     "strconv",
		},
	}

	expectedTypes := []TypeDecl{
		{
			Position: Position{Path: path, Line: 10},
			Name:     "ExaType",
			Type:     Slice{Elem: Map{Key: Int, Value: Struct{Fields: []StructField{{Name: "Err", Type: Error, Tag: "tag"}}}}},
		},
		{
			Position: Position{Path: path, Line: 14},
			Name:     "Alias",
			Alias:    Selector{Left: "io", Right: "ReadWriter"},
		},
		{
			Position: Position{Path: path, Line: 20},
			Name:     "Smoo",
			Type:     Int,
		},
	}

	expectedFuncs := []FuncDecl{
		{
			Position: Position{Path: path, Line: 16},
			Name:     "String",
			Recv:     "ExaType",
			Type:     Func{Results: []FuncResult{{Type: String}}},
		},
		{
			Position: Position{Path: path, Line: 30},
			Name:     "hello",
			Type:     Func{},
		},
	}

	if !reflect.DeepEqual(expectedImports, f.Imports) {
		t.Logf("want: %#v", expectedImports)
		t.Logf(" got: %#v", f.Imports)
		t.Fatalf("invalid imports")
	}

	if !reflect.DeepEqual(expectedTypes, f.Types) {
		t.Logf("want: %#v", expectedTypes)
		t.Logf(" got: %#v", f.Types)
		t.Fatalf("invalid types")
	}

	if !reflect.DeepEqual(expectedFuncs, f.Funcs) {
		t.Logf("want: %#v", expectedFuncs)
		t.Logf(" got: %#v", f.Funcs)
		t.Fatalf("invalid funcs")
	}

	expectedMethods := map[string]Func{
		"String": {Results: []FuncResult{{Type: String}}},
	}
	methods := f.GetMethods("ExaType")
	if !reflect.DeepEqual(expectedMethods, methods) {
		t.Logf("want: %#v", expectedMethods)
		t.Logf(" got: %#v", methods)
		t.Fatalf("invalid methods")
	}

	expectedGetTypes := map[string]Type{
		"ExaType": Slice{Elem: Map{Key: Int, Value: Struct{Fields: []StructField{{Name: "Err", Type: Error, Tag: "tag"}}}}},
		"Smoo":    Int,
	}
	types := f.GetTypes()
	if !reflect.DeepEqual(expectedGetTypes, types) {
		t.Logf("want: %#v", expectedGetTypes)
		t.Logf(" got: %#v", types)
		t.Fatalf("invalid methods")
	}
}

package gadget

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	String  = Ident("string")
	Bool    = Ident("bool")
	Byte    = Ident("byte")
	Rune    = Ident("rune")
	Uintptr = Ident("uintptr")

	Int   = Ident("int")
	Int8  = Ident("int8")
	Int16 = Ident("int16")
	Int32 = Ident("int32")
	Int64 = Ident("int64")

	Uint   = Ident("uint")
	Uint8  = Ident("uint8")
	Uint16 = Ident("uint16")
	Uint32 = Ident("uint32")
	Uint64 = Ident("uint64")

	Float32    = Ident("float32")
	Float64    = Ident("float64")
	Complex64  = Ident("complex64")
	Complex128 = Ident("complex128")

	EmptyStruct = Struct{}
	Bytes       = Slice{Elem: Byte}
	Error       = Ident("error")
)

// Type is one of the following concrete types:
// Ident, Pointer, Slice, Array, Map, Struct, Chan, Func, Interface, Selector
type Type interface {
	String() string
	isType()
}

type Ident string

func (i Ident) String() string {
	return string(i)
}

func (i Ident) isType() {}

type Selector struct {
	Left  Ident
	Right Ident
}

func (s Selector) String() string {
	return s.Left.String() + "." + s.Right.String()
}

func (s Selector) isType() {}

type Pointer struct {
	Elem Type
}

func (p Pointer) String() string {
	return "*" + p.Elem.String()
}

func (p Pointer) isType() {}

type Slice struct {
	Elem Type
}

func (s Slice) String() string {
	return "[]" + s.Elem.String()
}

func (s Slice) isType() {}

type Array struct {
	Elem Type
	Size int
}

func (a Array) String() string {
	return fmt.Sprintf("[%d]%s", a.Size, a.Elem.String())
}

func (a Array) isType() {}

type Map struct {
	Key   Type
	Value Type
}

func (m Map) String() string {
	return fmt.Sprintf("map[%s]%s", m.Key.String(), m.Value.String())
}

func (m Map) isType() {}

type Struct struct {
	Fields []StructField
}

func (s Struct) String() string {
	fields := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		fields[i] = field.String()
	}
	return fmt.Sprintf("struct{%s}", strings.Join(fields, "; "))
}

func (s Struct) isType() {}

type StructField struct {
	Name string
	Type Type
	Tag  string
}

func (f StructField) String() string {
	var full string
	if f.Name != "" {
		full = f.Name + " "
	}
	full += f.Type.String()
	if f.Tag != "" {
		full += " " + strconv.Quote(f.Tag)
	}
	return full
}

// ChanDir represents the direction of a channel type
type ChanDir int

const (
	SEND ChanDir = iota
	RECV
	BOTH
)

func (d ChanDir) String() string {
	switch d {
	case SEND:
		return "SEND"
	case RECV:
		return "RECV"
	case BOTH:
		return "BOTH"
	default:
		return "UNKNOWN"
	}
}

type Chan struct {
	Dir  ChanDir
	Elem Type
}

func (c Chan) String() string {
	switch c.Dir {
	case SEND:
		return fmt.Sprintf("chan<- %s", c.Elem.String())
	case RECV:
		return fmt.Sprintf("<-chan %s", c.Elem.String())
	case BOTH:
		return fmt.Sprintf("chan %s", c.Elem.String())
	}
	return fmt.Sprintf("unknown_chan_dir %s", c.Elem.String())
}

func (c Chan) isType() {}

type Func struct {
	Params  []FuncParam
	Results []FuncResult
}

func (f Func) String() string {
	return "func" + f.toPrototype()
}

func (f Func) isType() {}

func (f Func) toPrototype() string {
	var params []string
	for _, param := range f.Params {
		params = append(params, param.String())
	}
	full := fmt.Sprintf("(%s)", strings.Join(params, ", "))
	if len(f.Results) == 0 {
		return full
	}
	if len(f.Results) == 1 && f.Results[0].Name == "" {
		return full + " " + f.Results[0].Type.String()
	}
	var results []string
	for _, result := range f.Results {
		results = append(results, result.String())
	}
	return fmt.Sprintf("%s (%s)", full, strings.Join(results, ", "))
}

type FuncParam struct {
	Name string
	Type Type
}

func (f FuncParam) String() string {
	if f.Name == "" {
		return f.Type.String()
	}
	return f.Name + " " + f.Type.String()
}

type FuncResult struct {
	Name string
	Type Type
}

func (f FuncResult) String() string {
	if f.Name == "" {
		return f.Type.String()
	}
	return f.Name + " " + f.Type.String()
}

type Interface struct {
	Methods []InterfaceMethod
}

func (i Interface) String() string {
	fields := make([]string, len(i.Methods))
	for i, field := range i.Methods {
		fields[i] = field.String()
	}
	return fmt.Sprintf("interface {%s}", strings.Join(fields, "; "))
}

func (i Interface) isType() {}

type InterfaceMethod struct {
	Name string
	Type Func
}

func (f InterfaceMethod) String() string {
	return f.Name + f.Type.toPrototype()
}

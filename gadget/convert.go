package gadget

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func convertTypeSpec(spec ast.Expr) (Type, error) {
	switch t := spec.(type) {
	case *ast.Ident:
		return Ident(t.Name), nil
	case *ast.SelectorExpr:
		ident, ok := t.X.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("expected selector expression to start with identifier")
		}
		return Selector{
			Left:  Ident(ident.Name),
			Right: Ident(t.Sel.Name),
		}, nil
	case *ast.StarExpr:
		elem, err := convertTypeSpec(t.X)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pointer: %w", err)
		}
		return Pointer{
			Elem: elem,
		}, nil
	case *ast.ArrayType:
		elem, err := convertTypeSpec(t.Elt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert array element: %w", err)
		}
		if t.Len == nil {
			return Slice{Elem: elem}, nil
		}
		size, err := asIntLiteral(t.Len)
		if err != nil {
			return nil, fmt.Errorf("failed to convert array size: %w", err)
		}
		return Array{Elem: elem, Size: size}, nil
	case *ast.MapType:
		key, err := convertTypeSpec(t.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to convert map key: %w", err)
		}
		val, err := convertTypeSpec(t.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert map value: %w", err)
		}
		return Map{Key: key, Value: val}, nil
	case *ast.ChanType:
		elem, err := convertTypeSpec(t.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert channel element: %w", err)
		}
		dir := BOTH
		if t.Dir == ast.RECV {
			dir = RECV
		}
		if t.Dir == ast.SEND {
			dir = SEND
		}
		return Chan{
			Dir:  dir,
			Elem: elem,
		}, nil
	case *ast.StructType:
		var s Struct
		for fieldNum, field := range t.Fields.List {
			if field.Names == nil {
				t, err := convertTypeSpec(field.Type)
				if err != nil {
					return nil, fmt.Errorf("failed to convert struct field %d: %w", fieldNum+1, err)
				}
				var tag string
				if field.Tag != nil {
					tag, err = asStringLiteral(field.Tag)
					if err != nil {
						return nil, fmt.Errorf("failed to convert tag for struct field %d: %w", fieldNum+1, err)
					}
				}
				s.Fields = append(s.Fields, StructField{
					Type: t,
					Tag:  tag,
				})
				continue
			}
			for _, name := range field.Names {
				t, err := convertTypeSpec(field.Type)
				if err != nil {
					return nil, fmt.Errorf("failed to convert struct field %s: %w", name.Name, err)
				}
				var tag string
				if field.Tag != nil {
					tag, err = asStringLiteral(field.Tag)
					if err != nil {
						return nil, fmt.Errorf("failed to convert tag for struct field %s: %w", name.Name, err)
					}
				}
				s.Fields = append(s.Fields, StructField{
					Name: name.Name,
					Type: t,
					Tag:  tag,
				})
			}
		}
		return s, nil
	case *ast.FuncType:
		var params []FuncParam
		if t.Params != nil {
			for fieldNum, field := range t.Params.List {
				if field.Names == nil {
					t, err := convertTypeSpec(field.Type)
					if err != nil {
						return nil, fmt.Errorf("failed to convert function parameter %d: %w", fieldNum+1, err)
					}
					params = append(params, FuncParam{
						Type: t,
					})
					continue
				}
				for _, name := range field.Names {
					t, err := convertTypeSpec(field.Type)
					if err != nil {
						return nil, fmt.Errorf("failed to convert function parameter %s: %w", name.Name, err)
					}
					params = append(params, FuncParam{
						Name: name.Name,
						Type: t,
					})
				}
			}
		}
		var results []FuncResult
		if t.Results != nil {
			for fieldNum, field := range t.Results.List {
				if field.Names == nil {
					t, err := convertTypeSpec(field.Type)
					if err != nil {
						return nil, fmt.Errorf("failed to convert function return value %d: %w", fieldNum+1, err)
					}
					results = append(results, FuncResult{
						Type: t,
					})
					continue
				}
				for _, name := range field.Names {
					t, err := convertTypeSpec(field.Type)
					if err != nil {
						return nil, fmt.Errorf("failed to convert function return value %s: %w", name.Name, err)
					}
					results = append(results, FuncResult{
						Name: name.Name,
						Type: t,
					})
				}
			}
		}
		return Func{
			Params:  params,
			Results: results,
		}, nil
	case *ast.InterfaceType:
		var i Interface
		if t.Methods != nil {
			for _, field := range t.Methods.List {
				for _, name := range field.Names {
					t, err := convertTypeSpec(field.Type)
					if err != nil {
						return nil, fmt.Errorf("failed to convert function parameter %s: %w", name.Name, err)
					}
					f, ok := t.(Func)
					if !ok {
						return nil, fmt.Errorf("interface method %s was somehow not a function type: %w", name.Name, err)
					}
					i.Methods = append(i.Methods, InterfaceMethod{
						Name: name.Name,
						Type: f,
					})
				}
			}
		}
		return i, nil
	}
	return nil, fmt.Errorf("unknown kind of type spec %#v", spec)
}

func asIntLiteral(expr ast.Expr) (int, error) {
	t, ok := expr.(*ast.BasicLit)
	if !ok {
		return 0, fmt.Errorf("expression is not a basic literal")
	}
	if t.Kind != token.INT {
		return 0, fmt.Errorf("expression is not a basic literal integer")
	}
	return strconv.Atoi(t.Value)
}

func asStringLiteral(expr ast.Expr) (string, error) {
	t, ok := expr.(*ast.BasicLit)
	if !ok {
		return "", fmt.Errorf("expression is not a basic literal")
	}
	if t.Kind != token.STRING {
		return "", fmt.Errorf("expression is not a basic literal string")
	}
	return strconv.Unquote(t.Value)
}

type funcVisitor func(ast.Node) ast.Visitor

func (f funcVisitor) Visit(node ast.Node) ast.Visitor {
	return f(node)
}

func findTypeAlias(fileSet *ast.File, pos token.Pos) (Type, error) {
	done := false
	var typ ast.Expr
	var err error

	var search funcVisitor
	search = funcVisitor(func(node ast.Node) ast.Visitor {
		if done {
			return nil
		}
		if node == nil {
			return search
		}
		if node.Pos() > pos {
			done = true
			switch concrete := node.(type) {
			case *ast.Ident:
				typ = concrete
				return nil
			case *ast.SelectorExpr:
				typ = concrete
				return nil
			default:
				err = fmt.Errorf("type alias is not an identifier or selector")
				return nil
			}
		}
		return search
	})
	ast.Walk(search, fileSet)
	if err != nil {
		return nil, fmt.Errorf("failed to find alias: %w", err)
	}
	if typ == nil {
		return nil, fmt.Errorf("failed to find alias")
	}
	converted, err := convertTypeSpec(typ)
	if err != nil {
		return nil, fmt.Errorf("failed to convert identifier in alias: %w", err)
	}
	return converted, nil
}

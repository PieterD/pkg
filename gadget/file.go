package gadget

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

// File contains all the information we have about a parsed Go file.
type File struct {
	Path      string       // The path of the Go file this File represents.
	Imports   []ImportDecl // The imports contained within the Go file.
	Types     []TypeDecl   // The Type declarations contained within the Go file.
	Funcs     []FuncDecl   // The Function declarations contained within the Go file.
	HasErrors bool         // HasErrors is true if there was an invalid declaration was found.
}

// Position represents a file:line location.
type Position struct {
	Path string
	Line int
}

func (pos Position) String() string {
	return fmt.Sprintf("%s:%d", pos.Path, pos.Line)
}

type ImportDecl struct {
	Position
	Name string // The identifier of the import decl (can be empty "." or "_")
	Path string // The import path.
}

type TypeDecl struct {
	Position
	Name  string // The type name.
	Type  Type   // The actual type definition. May be empty if the type declaration is an alias.
	Alias Type   // The alias this declaration references. May be nil if the type declaration is not an alias. Is either Ident or Selector.
}

type FuncDecl struct {
	Position
	Name string // The function name.
	Recv string // The receiver type identifier.
	Type Func   // The function type.
}

// NewFile parses a Go file.
// If reader is nil, the file at path is opened.
// Otherwise, reader is taken to be the contents of the file.
func NewFile(path string, reader io.Reader) (*File, error) {
	if reader == nil {
		h, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open file '%s': %w", path, err)
		}
		defer h.Close()
		reader = h
	}
	f := &File{
		Path: path,
	}
	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, path, reader, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%s': %w", path, err)
	}

	for _, decl := range parsedFile.Decls {
		pos := Position{
			Path: path,
			Line: fileSet.Position(decl.Pos()).Line,
		}
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				pos := Position{
					Path: path,
					Line: fileSet.Position(spec.Pos()).Line,
				}
				switch decl.Tok {
				case token.TYPE:
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						return nil, fmt.Errorf("%s: expected *ast.TypeSpec, got %T", pos, typeSpec)
					}
					name := typeSpec.Name.Name
					var (
						alias Type
						typ   Type
					)
					if typeSpec.Assign.IsValid() {
						alias, err = findTypeAlias(parsedFile, typeSpec.Assign)
						if err != nil {
							return nil, fmt.Errorf("%s: failed to parse alias for type '%s': %w", pos, name, err)
						}
					} else {
						typ, err = convertTypeSpec(typeSpec.Type)
						if err != nil {
							return nil, fmt.Errorf("%s: failed to convert type %s: %w", pos, name, err)
						}
					}
					f.Types = append(f.Types, TypeDecl{
						Position: pos,
						Name:     name,
						Type:     typ,
						Alias:    alias,
					})
				case token.VAR:
					varSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						return nil, fmt.Errorf("%s: expected *ast.ValueSpec, got %T", pos, varSpec)
					}
				case token.CONST:
					varSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						return nil, fmt.Errorf("%s: expected *ast.ValueSpec, got %T", pos, varSpec)
					}
				case token.IMPORT:
					importSpec, ok := spec.(*ast.ImportSpec)
					if !ok {
						return nil, fmt.Errorf("%s: expected *ast.ImportSpec, got %T", pos, importSpec)
					}
					impName := ""
					if importSpec.Name != nil {
						impName = importSpec.Name.Name
					}
					impPath, err := asStringLiteral(importSpec.Path)
					if err != nil {
						return nil, fmt.Errorf("%s: failed to parse import path: %w", pos, err)
					}
					f.Imports = append(f.Imports, ImportDecl{
						Position: pos,
						Name:     impName,
						Path:     impPath,
					})
				}
			}
		case *ast.FuncDecl:
			recv := ""
			if decl.Recv != nil && len(decl.Recv.List) != 0 {
				if len(decl.Recv.List) > 1 {
					return nil, fmt.Errorf("%s: multiple method receivers", pos)
				}
				typ, err := convertTypeSpec(decl.Recv.List[0].Type)
				if err != nil {
					return nil, fmt.Errorf("%s: failed to convert method receiver type: %w", pos, err)
				}
				ptr, ok := typ.(Pointer)
				if ok {
					typ = ptr.Elem
				}
				id, ok := typ.(Ident)
				if !ok {
					return nil, fmt.Errorf("%s: method receiver type is not Identifier or *Identifier", pos)
				}
				recv = id.String()
			}
			typ, err := convertTypeSpec(decl.Type)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to convert function type: %w", pos, err)
			}
			t, ok := typ.(Func)
			if !ok {
				return nil, fmt.Errorf("%s: function declaration type is somehow not a function type", pos)
			}
			f.Funcs = append(f.Funcs, FuncDecl{
				Position: pos,
				Name:     decl.Name.Name,
				Recv:     recv,
				Type:     t,
			})
		case *ast.BadDecl:
			f.HasErrors = true
		}
	}

	return f, nil
}

// GetMethods fetches the methods belonging to the given type identifier.
func (f *File) GetMethods(typeName string) map[string]Func {
	decls := make(map[string]Func)
	for _, fun := range f.Funcs {
		if fun.Recv == typeName {
			decls[fun.Name] = fun.Type
		}
	}
	if len(decls) == 0 {
		return nil
	}
	return decls
}

// GetTypes fetches all non-alias types.
func (f *File) GetTypes() map[string]Type {
	decls := make(map[string]Type)
	for _, decl := range f.Types {
		if decl.Type == nil {
			continue
		}
		decls[decl.Name] = decl.Type
	}
	if len(decls) == 0 {
		return nil
	}
	return decls
}

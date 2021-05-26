package go_comments_unmarshaler

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
)

func UnmarshalPackage(pathToPkg string, result interface{}) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(result)}
	}
	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("only structs supported, but got %q", rv.Elem().Elem().Kind())
	}

	return unmarshalPackage(pathToPkg, newVisitor(rv.Elem()))
}

func UnmarshalModule(pathToRoot string, result interface{}) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(result)}
	}
	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("only structs supported, but got %q", rv.Elem().Elem().Kind())
	}
	visitor := newVisitor(rv.Elem())
	return filepath.WalkDir(pathToRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(pathToRoot, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			relPath = ""
		} else {
			relPath += "/"
		}
		return unmarshalPackage(path, visitor.withPath(relPath))
	})
}

func unmarshalPackage(pathToPkg string, visitor *pkgVisitor) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pathToPkg, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Walk(visitor, file)
		}
	}
	return nil
}

type pkgVisitor struct {
	currentPath string
	target      reflect.Value
	written     map[string]bool
}

func (p *pkgVisitor) withPath(currentPath string) *pkgVisitor {
	return &pkgVisitor{
		target:      p.target,
		currentPath: currentPath,
		written:     p.written,
	}
}

func newVisitor(rv reflect.Value) *pkgVisitor {
	return &pkgVisitor{
		target:      rv,
		written:     map[string]bool{},
		currentPath: "",
	}
}

func (p *pkgVisitor) Visit(node ast.Node) ast.Visitor {
	if !p.wantVisitChildren(node) {
		return nil
	}
	if _, ok := node.(*ast.File); ok {
		return p
	}
	p.collect(node, p.currentPath)
	return p
}

func (p *pkgVisitor) wantVisitChildren(node ast.Node) bool {
	switch node.(type) {
	case *ast.GenDecl, *ast.FuncDecl, *ast.File:
		return true
	default:
		return false
	}
}

func (p *pkgVisitor) collect(n ast.Node, currentPath string) {
	switch node := n.(type) {
	case *ast.FuncDecl:
		currentPath = prepareFuncPath(node, currentPath)
		findValueByPath(p.target, currentPath, node.Doc.Text())
	case *ast.GenDecl:
		switch node.Tok {
		case token.TYPE:
			n := len(node.Specs)
			if n == 0 {
				return
			}
			for i := range node.Specs {
				ts, ok := node.Specs[i].(*ast.TypeSpec)
				if !ok {
					// impossible
					continue
				}
				path := currentPath + ts.Name.Name
				var text string
				// only one declaration in format: `type X olala`
				if n == 1 {
					text = node.Doc.Text()
				} else {
					// multiple type declaration in one parentheses
					text = ts.Doc.Text()
				}
				findValueByPath(p.target, path, text)
			}
		default:
			// other not supported yet
			return
		}
	}
}

func findValueByPath(root reflect.Value, path, text string) bool {
	switch root.Kind() {
	case reflect.String:
		root.SetString(text)
		return true
	case reflect.Struct:
		rtype := root.Type()
		for i, n := 0, root.NumField(); i < n; i++ {
			ftype := rtype.Field(i)
			tag := ftype.Tag.Get("comment")
			if tag == "" {
				continue
			}
			field := root.Field(i)
			if field.Kind() != reflect.String {
				if !strings.HasPrefix(path, tag) && tag != "*" {
					continue
				}
				path = strings.TrimPrefix(path, tag)
				path = strings.TrimPrefix(path, "/")
				ok := findValueByPath(field, path, text)
				if ok {
					return true
				}
				continue
			}
			if tag == path {
				findValueByPath(field, path, text)
				return true
			}
		}
		return false
	case reflect.Map:
		rtype := root.Type()
		if rtype.Key().Kind() != reflect.String {
			panic("maps with not string keys is not supported")
		}
		if rtype.Elem().Kind() != reflect.Struct {
			panic("maps with not struct values is not supported")
		}
		if root.IsNil() {
			root.Set(reflect.MakeMap(rtype))
		}
		splitted := strings.SplitN(path, "/", 2)
		if len(splitted) < 2 {
			return false
		}
		key := reflect.New(rtype.Key()).Elem()
		key.SetString(splitted[0])

		found := root.MapIndex(key)
		mval := reflect.New(rtype.Elem()).Elem()
		if found.IsValid() {
			mval.Set(found)
		}
		ok := findValueByPath(mval, splitted[1], text)
		if !ok {
			return false
		}
		root.SetMapIndex(key, mval)
		return true
	default:
		return false
	}
}

func prepareFuncPath(node *ast.FuncDecl, currentPath string) string {
	// parse method
	if node.Recv != nil && len(node.Recv.List) > 0 {
		recType := node.Recv.List[0].Type
		// Here can be Type or *Type
		if star, ok := recType.(*ast.StarExpr); ok {
			recType = star.X
		}
		switch rec := recType.(type) {
		case *ast.Ident:
			currentPath += rec.Name + "."
		default:
			panic(fmt.Sprintf("cant parse %T in method receiver", recType))
		}
	}

	currentPath += node.Name.Name
	return currentPath
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

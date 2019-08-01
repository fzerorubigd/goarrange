package goarrange

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

func getName(id ast.Expr) string {
	switch t := id.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return getName(t.X)
	default:
		panic(fmt.Sprintf("unsupported type %T", t))
	}
}

type decl struct {
	ast.Decl
	receiver string
	name     string
}

type declTree []decl

func (dt declTree) Len() int {
	return len(dt)
}

func (dt declTree) Less(i, j int) bool {
	if dt[i].receiver != "" && dt[j].receiver == "" {
		return true
	}

	if dt[i].receiver == "" && dt[j].receiver != "" {
		return false
	}

	if dt[i].receiver != dt[j].receiver {
		return strings.Compare(dt[i].receiver, dt[j].receiver) < 0
	}

	return strings.Compare(dt[i].name, dt[j].name) < 0
}

func (dt declTree) Swap(i, j int) {
	dt[i], dt[j] = dt[j], dt[i]
}

// Arrange the go source file in alphabetical order
func Arrange(src []byte) ([]byte, error) {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var tree declTree
	for i := range f.Decls {
		switch t := f.Decls[i].(type) {
		case *ast.FuncDecl:
			rec := ""
			if t.Recv != nil {
				rec = getName(t.Recv.List[0].Type)
			}
			tree = append(tree, decl{
				name:     getName(t.Name),
				receiver: rec,
				Decl:     t,
			})
		case *ast.GenDecl:
			fmt.Println(t.Tok)
		default:
			return nil, fmt.Errorf("type %T is not supported", t)
		}
	}

	sort.Sort(&tree)

	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buf, "package %s\n", getName(f.Name))

	for i := range tree {
		_, _ = fmt.Fprint(buf, "\n")
		if err := format.Node(buf, fSet, tree[i].Decl); err != nil {
			return nil, err
		}
		_, _ = fmt.Fprint(buf, "\n")
	}

	return buf.Bytes(), nil
}

package finder

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Walker struct {
	state
	Name    string
	Package string
}

type state struct {
	pack string
	fn   string
}

func (w *Walker) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.AssignStmt:
		if d.Tok != token.DEFINE {
			return w
		}

		for _, name := range d.Rhs {
			if ident, ok := name.(*ast.Ident); ok {
				if ident.Name == w.Name {
					// find the assignment
				}
			}
		}
	case *ast.File:
		w.pack = d.Name.Name
		for _, decl := range d.Decls {
			if fdecl, ok := decl.(*ast.FuncDecl); ok {
				for _, field := range fdecl.Type.Params.List {
					if star, ok := field.Type.(*ast.StarExpr); ok {
						if s, ok := star.X.(*ast.SelectorExpr); ok {
							ident := s.X.(*ast.Ident).Name
							if ident == w.Package && s.Sel.Name == w.Name {
								fmt.Printf("(pkg: %s, fn: %s) as fn param: *%s.%s\n", w.pack, fdecl.Name, ident, s.Sel.Name)
							}
							continue
						}
						if ident, ok := star.X.(*ast.Ident); ok && ident.Name == w.Name {
							fmt.Printf("(pkg: %s, fn: %s) as fn param: *%s\n", w.pack, fdecl.Name, ident.Name)
						}
					}
				}
			}
		}
	}
	return w
}

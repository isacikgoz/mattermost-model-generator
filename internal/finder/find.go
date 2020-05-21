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
		// find in function declarations
		for _, decl := range d.Decls {
			fdecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			var rcv string
			if fdecl.Recv != nil {
				for _, field := range fdecl.Recv.List {
					switch f := field.Type.(type) {
					case *ast.Ident:
						rcv = f.Name + "."
					case *ast.StarExpr:
						rcv = f.X.(*ast.Ident).Name + "."
					}
				}
			}
			var varname, typStr string
			for _, field := range fdecl.Type.Params.List {
				// if we pass the pointer proceed
				star, ok := field.Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				var found bool
				// same package?
				if ident, ok := star.X.(*ast.Ident); ok && ident.Name == w.Name {
					found = true
					typStr = fmt.Sprintf("*%s", ident.Name)
				}
				// other package
				if s, ok := star.X.(*ast.SelectorExpr); ok {
					ident := s.X.(*ast.Ident).Name
					if ident == w.Package && s.Sel.Name == w.Name {
						found = true
						typStr = fmt.Sprintf("*%s.%s", ident, s.Sel.Name)
					}
				}

				if !found {
					continue
				}

				for _, name := range field.Names {
					varname = name.Name // TODO: what if multiple?
				}
			}
			// we have a variable, find if something is being assigned in the function block
			if varname != "" {
				if found := findAssignmentsInBlock(varname, fdecl.Body); found {
					fmt.Printf("(pkg: %s, fn: %s%s) as fn param: %s %s\n", w.pack, rcv, fdecl.Name, varname, typStr)
				}
			}
		}
	}
	return w
}

func findAssignmentsInBlock(varname string, block *ast.BlockStmt) bool {
	var assigned bool
	for _, a := range block.List {
		// look to assignments
		as, ok := a.(*ast.AssignStmt)
		if !ok {
			continue
		}
		// check if we are assigning something to our variable
		// in the left hand side
		for _, expression := range as.Lhs {
			selector, ok := expression.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			ident, ok := selector.X.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name != varname {
				continue
			}
			assigned = true
			// something is assigned to our value
			fmt.Printf("%s.%s (assignment)\n", varname, selector.Sel.Name)
		}
	}
	return assigned
}

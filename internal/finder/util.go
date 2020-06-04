package finder

import (
	"go/ast"
)

// eg. vartype: Channel, pkg: model
// returns true if there is a mutation for the given object
func findMutationForFunctionArgument(vartype, varpackage, functionpackage string, fdecl *ast.FuncDecl) (map[*ast.Field][]*ast.AssignStmt, bool) {
	var assigned bool
	mutatedFields := make(map[*ast.Field][]*ast.AssignStmt)

	for _, field := range fdecl.Type.Params.List {
		// if we pass the pointer proceed
		star, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		var found bool
		// same package?
		if ident, ok := star.X.(*ast.Ident); ok && (ident.Name == vartype) && (varpackage == functionpackage) {
			found = true
		}
		// other package
		if s, ok := star.X.(*ast.SelectorExpr); ok && !found {
			ident := s.X.(*ast.Ident).Name
			found = ident == varpackage && s.Sel.Name == vartype
		}

		if !found {
			continue
		}

		for _, name := range field.Names {
			if assignements, found := findAssignmentsInBlock(name.Name, fdecl.Body); found {
				assigned = true
				mutatedFields[field] = assignements
			}
		}
	}

	return mutatedFields, assigned
}

func findAssignmentsInBlock(varname string, block *ast.BlockStmt) ([]*ast.AssignStmt, bool) {
	var assignments []*ast.AssignStmt
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
			assignments = append(assignments, as)
			assigned = true
		}
	}
	return assignments, assigned
}

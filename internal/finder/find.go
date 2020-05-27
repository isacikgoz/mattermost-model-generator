package finder

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type Walker struct {
	state
	Name    string
	Package string
}

type state struct {
	pack       string
	paramLists map[*ast.FieldList][]*ast.Field // map of field list to specific fields to rename
	funcLists  map[*ast.FuncDecl]interface{}   // map of field list to specific fields to rename
}

func renameModelParameter(field *ast.Field) *ast.Field {
	return &ast.Field{
		Comment: field.Comment,
		Doc:     field.Doc,
		Tag:     field.Tag,
		Type:    field.Type,
		Names: []*ast.Ident{
			{
				Name: "_" + field.Names[0].Name,
			},
		},
	}
}

func addCloneToBlock(block *ast.BlockStmt, fields []*ast.Field) *ast.BlockStmt {
	assignments := []ast.Stmt{}
	for _, field := range fields {
		assignments = append(assignments, &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: field.Names[0].Name,
				},
			},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "_" + field.Names[0].Name},
						Sel: &ast.Ident{Name: "Clone"},
					},
				},
			},
		})
	}
	return &ast.BlockStmt{
		Lbrace: block.Lbrace,
		Rbrace: block.Rbrace,
		List:   append(assignments, block.List...),
	}
}

func (w *Walker) pre(c *astutil.Cursor) bool {
	if parentFunc, ok := c.Parent().(*ast.FuncDecl); ok && w.funcLists[parentFunc] == parentFunc {
		if stmt, ok := c.Node().(*ast.BlockStmt); ok {
			c.Replace(addCloneToBlock(stmt, w.paramLists[parentFunc.Type.Params]))
			return true
		}
	}
	if parentList, ok := c.Parent().(*ast.FieldList); ok && w.paramLists[parentList] != nil {
		for _, ff := range w.paramLists[parentList] {
			if ff == c.Node() {
				c.Replace(renameModelParameter(ff))
				return true
			}
		}
	}
	if fdecl, ok := c.Node().(*ast.FuncDecl); ok {
		fields, ok := findMutationForFunctionArgument(w.Name, w.Package, w.pack, fdecl)
		if !ok {
			return true
		}
		fmt.Printf("pkg: %s, fn: %s\n", w.pack, fdecl.Name)

		w.paramLists[fdecl.Type.Params] = fields
		w.funcLists[fdecl] = fdecl
	}
	return true
}

func (w *Walker) post(c *astutil.Cursor) bool {
	return true
}

func (w *Walker) Process(root ast.Node) ast.Node {
	w.paramLists = make(map[*ast.FieldList][]*ast.Field)
	w.funcLists = make(map[*ast.FuncDecl]interface{})
	w.pack = root.(*ast.File).Name.Name
	return astutil.Apply(root, w.pre, w.post)
}

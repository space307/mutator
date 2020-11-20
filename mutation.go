package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const (
	mFunctionPrefix = "MutateTo"
)

type (
	// MParser mutation files parser
	MParser struct {
		m MFuncList
	}

	// MFuncList contain list of mutation functions
	MFuncList []MFunc

	// MFunc mutation function description
	MFunc struct {
		Name     string
		Receiver MStruct
		Result   MStruct
	}

	// MStruct mutation struct
	MStruct struct {
		Pkg  string
		Name string
	}

	mVisitor struct {
		p    *MParser
		file *ast.File
	}
)

// NewMParser constructor
func NewMParser() *MParser {
	return &MParser{
		m: MFuncList{},
	}
}

// Parse mutator functions from file
func (p *MParser) Parse(filename string) (MFuncList, error) {
	p.m = p.m[:0]
	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		return nil, err
	}
	ast.Walk(&mVisitor{p: p}, f)
	return p.m, nil
}

// Visit implementation
func (v *mVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Package:
		return v
	case *ast.File:
		v.file = n
		return v
	case *ast.FuncDecl:
		if !strings.HasPrefix(mFunctionPrefix, n.Name.Name) {
			return nil
		}
		recv := receiverStruct(v.file.Name.Name, n)
		if recv == nil {
			return nil
		}
		res := resultStruct(n)
		if res == nil {
			return nil
		}
		if n.Name.Name != mFunctionPrefix+res.Name {
			return nil
		}
		v.p.m = append(v.p.m, MFunc{
			Name:     n.Name.Name,
			Receiver: *recv,
			Result:   *res,
		})
	}

	return nil
}

func receiverStruct(pkg string, n *ast.FuncDecl) *MStruct {
	if len(n.Recv.List) != 1 {
		return nil
	}
	rType, ok := n.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return nil
	}
	ident, ok := rType.X.(*ast.Ident)
	if !ok {
		return nil
	}
	return &MStruct{
		Pkg:  pkg,
		Name: ident.Name,
	}
}

func resultStruct(n *ast.FuncDecl) *MStruct {
	if n.Type.Results == nil || len(n.Type.Results.List) != 1 {
		return nil
	}
	selector, ok := n.Type.Results.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return nil
	}
	return &MStruct{
		Pkg:  ident.Name,
		Name: selector.Sel.Name,
	}
}

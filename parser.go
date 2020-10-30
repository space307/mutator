package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const marker = "mutagento"

type Struct struct {
	Path string
	Name string
}

type Parser struct {
	PkgPath     string
	PkgName     string
	StructNames []Struct
	tmpFile     string

	All         map[string]*ast.StructType
	Connections map[Struct]Struct
}

type visitor struct {
	*Parser

	tempPkgName string
	name        string
}

func (v *visitor) Visit(n ast.Node) (w ast.Visitor) {
	switch n := n.(type) {
	case *ast.Package:
		return v
	case *ast.File:
		// тут запоминаем путь до пакета в котором структура, которую парсим
		v.PkgName = n.Name.String()
		v.tempPkgName = filepath.Dir(n.Name.String())
		return v
	case *ast.TypeSpec:

		structType, ok := n.Type.(*ast.StructType)
		if !ok {
			// если не структура ищем дальше
			return nil
		}

		// кладём во все структуры
		v.All[strings.TrimSpace(n.Name.String())] = structType

		//проверяем коммент на маркер, если он найдет, то запоминаем пару
		if st := checkConnect(n.Doc.Text()); st != nil {
			v.Connections[Struct{
				Path: v.tmpFile,
				Name: strings.TrimSpace(n.Name.String()),
			}] = *st
		}

		return v
	}
	return nil
}

func (p *Parser) Parse(fname string) error {
	err := filepath.Walk(fname, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		p.tmpFile = filepath.Dir(path)
		//парсим структуры из файла
		ast.Walk(&visitor{Parser: p}, f)

		return nil
	})

	return err
}

func checkConnect(comment string) *Struct {
	if strings.Contains(comment, marker) {
		t := strings.Split(comment, " ")
		for i, l := range t {
			if l == marker {
				if len(t) > i+2 {
					return &Struct{
						Path: t[i+1],
						Name: strings.TrimSpace(t[i+2]),
					}
				}
			}
		}
	}
	return nil
}

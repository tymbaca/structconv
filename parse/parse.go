package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type StructInfo struct {
	Name       string
	TypeIdent  *ast.Ident
	TypeSpec   *ast.TypeSpec
	StructType *ast.StructType
}

func Parse(fileName string) []StructInfo {
	// The name of the file we want to parse
	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var specs []*ast.TypeSpec
	// fmt.Printf("%#v\n", node.Decls[len(node.Decls)])
	for _, decl := range node.Decls {
		curSpecs, ok := getStructTypeSpecFromDecl(decl)
		if !ok {
			continue
		}
		specs = append(specs, curSpecs...)
	}

	return makeStructInfos(specs)
}

func makeStructInfos(specs []*ast.TypeSpec) []StructInfo {
	var infos []StructInfo
	for _, t := range specs {
		info := StructInfo{}
		info.Name = t.Name.String()
		info.TypeIdent = ast.NewIdent(t.Name.String())
		info.TypeSpec = t

		s, ok := t.Type.(*ast.StructType)
		if !ok {
			panic("got not struct in makeStructInfos")
		}
		info.StructType = s

		infos = append(infos, info)
	}

	return infos
}

type (
	A int
	B struct {
		//
	}
)

func getStructTypeSpecFromDecl(d ast.Decl) ([]*ast.TypeSpec, bool) {
	g, ok := d.(*ast.GenDecl)
	if !ok {
		return nil, false
	}

	var tSpecs []*ast.TypeSpec
	for _, spec := range g.Specs {
		t, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		_, ok = t.Type.(*ast.StructType)
		if !ok {
			continue
		}

		tSpecs = append(tSpecs, t)
	}

	if len(tSpecs) == 0 {
		return nil, false
	}

	return tSpecs, true
}

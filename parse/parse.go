package parse

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
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

	return makeStructInfos(specs...)
}

func makeStructInfos(specs ...*ast.TypeSpec) []StructInfo {
	var infos []StructInfo
	for _, t := range specs {
		info := StructInfo{}
		info.Name = t.Name.String()
		info.TypeIdent = ast.NewIdent(t.Name.String())
		info.TypeSpec = t

		s, ok := t.Type.(*ast.StructType)
		if !ok {
			continue
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

func ParseFile(path string) (*ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func FindDeclByComment(node *ast.File, target string) (*ast.FuncDecl, bool) {
	for _, d := range node.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if DocContains(fd.Doc, target) {
			return fd, true
		}
	}

	return nil, false
}

func DocContains(doc *ast.CommentGroup, target string) bool {
	for _, c := range doc.List {
		if strings.Contains(c.Text, target) {
			return true
		}
	}

	return false
}

// WARN: does't support non-struct types
func GitParamsAndResults(fd *ast.FuncDecl) (StructInfo, StructInfo, error) {
	if len(fd.Type.Params.List) != 1 || len(fd.Type.Results.List) != 1 {
		return StructInfo{}, StructInfo{}, errors.New("param or result fields are no count 1")
	}

	paramField := fd.Type.Params.List[0]
	// resultField := fd.Type.Results.List[0]

	paramTypeIdent, ok := paramField.Type.(*ast.Ident)
	if !ok {
		panic("can't convert param field to ident")
	}

	log.Println(paramTypeIdent.Name)

	// _, ok := paramField.Type.(*ast.StructType)
	// if !ok {
	// 	panic("can't convert param field to struct type")
	// }
	// ast.G

	return StructInfo{}, StructInfo{}, nil
}

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

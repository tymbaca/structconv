package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"

	"github.com/tymbaca/structconv/parse"
)

var (
	_srcIdent = ast.NewIdent("src")
	_dstIdent = ast.NewIdent("dst")
)

type GenInfo struct {
	Src parse.StructInfo
	Dst parse.StructInfo
	// map - which dst field corresponds to which src field
	// it is in this order because ... fuck that, i will explain later
	FieldMatch map[uint]uint
}

func Generate(info GenInfo, outputPath string) {
	outFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, outputPath, nil, parser.SkipObjectResolution)
	node.Name = ast.NewIdent("main")

	converterFuncNode := generateConvertor(info)

	node.Decls = append(node.Decls, converterFuncNode)

	printer.Fprint(outFile, fset, node)
}

func GenerateInPlace(targetPath string) {
	genInfo, file := GenInfoAndNodeFromFile(targetPath)
	_, _ = genInfo, file

	// converterFuncNode := generateConvertor(info)

	// node.Decls = append(node.Decls, converterFuncNode)

	// printer.Fprint(outFile, fset, node)
}

func GenInfoAndNodeFromFile(path string) (GenInfo, *ast.File) {
	node, err := parse.ParseFile(path)
	if err != nil {
		panic(err)
	}

	fd, ok := parse.FindDeclByComment(node, "structconv")
	if !ok {
		panic("can'd find 'structconv' func decl")
	}

	src, dst, err := parse.GitParamsAndResults(fd)
	if err != nil {
		panic(err)
	}

	info := GenInfo{
		Src: src,
		Dst: dst,
	}

	return info, node
}

func generateConvertor(info GenInfo) *ast.FuncDecl {
	fd := &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("convert%sTo%s", info.Src.Name, info.Dst.Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{_srcIdent},
						Type:  info.Src.TypeIdent,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						// Names: []*ast.Ident{ast.NewIdent("dst")},
						Type: info.Dst.TypeIdent,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: info.Dst.TypeIdent,
							Elts: generateFieldsKeyVals(info),
						},
					},
				},
			},
		},
	}

	return fd
}

func generateFieldsKeyVals(info GenInfo) []ast.Expr {
	if len(info.FieldMatch) > info.Dst.StructType.Fields.NumFields() {
		panic("provided key-value count is more then Dst struct field count")
	}

	var kvs []ast.Expr
	for dIdx, sIdx := range info.FieldMatch {
		if sIdx >= uint(info.Src.StructType.Fields.NumFields()) || dIdx >= uint(info.Dst.StructType.Fields.NumFields()) {
			panic("provided field index is out of range")
		}

		sF := info.Src.StructType.Fields.List[sIdx]
		dF := info.Dst.StructType.Fields.List[dIdx]
		// WARN: fill panic if type is interface of struct or anithing else
		sFIdent := sF.Type.(*ast.Ident)
		dFIdent := dF.Type.(*ast.Ident)

		if sFIdent.Name != dFIdent.Name {
			log.Panicf("field type are different: %v %v - %v %v", sF.Names[0], sF.Type, dF.Names[0], dF.Type)
		}

		kv := &ast.KeyValueExpr{
			Key: dF.Names[0],
			Value: &ast.SelectorExpr{
				X:   _srcIdent,
				Sel: sF.Names[0],
			},
		}
		kvs = append(kvs, kv)
	}

	return kvs
}

package main

import (
	"github.com/tymbaca/structconv/generate"
)

func main() {
	srcPath := "source.go"
	// structInfos := parse.Parse(srcPath)
	// outPath := "output.go"
	// genIngo := generate.GenInfo{
	// 	Src: structInfos[0],
	// 	Dst: structInfos[2],
	// 	FieldMatch: map[uint]uint{
	// 		0: 0,
	// 		1: 0,
	// 	},
	// }
	generate.GenerateInPlace(srcPath)
}

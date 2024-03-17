package main

type Src struct {
	Name   string
	Nested *Nested
	Age    int
}

type Nested struct {
}

type Dst struct {
	Years int
	Who   string
}

// structconv
func convertSrcToDst(src Src) Dst

package analyze

import "github.com/arrietty-lang/arrtty/preprocess/parse"

type Semantics struct {
	KnownValues      map[string]map[int]map[string][]*parse.DataType
	KnownFunctions   map[string]*FnDataType
	OutsideValues    []*parse.Node
	OutsideFunctions []*parse.Node
	Tree             []*parse.Node
}

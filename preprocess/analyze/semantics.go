package analyze

import "github.com/arrietty-lang/arrtty/preprocess/parse"

type Semantics struct {
	KnownValues    map[string]map[string][]*parse.DataType
	KnownFunctions map[string][]*parse.DataType
	Tree           *parse.Node
}

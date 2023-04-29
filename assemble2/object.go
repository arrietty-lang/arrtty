package assemble2

import "github.com/arrietty-lang/arrtty/preprocess/analyze"

type Object struct {
	Identifier    string
	SemanticsNode *analyze.Semantics
}

package analyze

import "github.com/arrietty-lang/arrtty/preprocess/parse"

type FnDataType struct {
	Params  []*parse.DataType
	Returns []*parse.DataType
}

package parse

type FuncDefField struct {
	Identifier *Node
	Parameters *Node
	Returns    *Node
	Body       *Node
}

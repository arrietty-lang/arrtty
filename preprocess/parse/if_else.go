package parse

type IfElseField struct {
	UseElse   bool
	Cond      *Node
	IfBlock   *Node
	ElseBlock *Node
}

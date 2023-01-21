package parse

import (
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
)

type Node struct {
	Kind NodeKind
	Pos  *tokenize.Position

	ImportField *ImportField

	DataTypeField *DataTypeField
	IdentField    *IdentField
	VarDeclField  *VarDeclField
	AssignField   *AssignField
	CommentField  *CommentField
	FuncDefField  *FuncDefField
	BlockField    *BlockField
	ReturnField   *ReturnField
	IfElseField   *IfElseField
	ForField      *ForField
}

func NewNode(kind NodeKind, pos *tokenize.Position) *Node {
	return &Node{
		Kind: kind,
		Pos:  pos,
	}
}

func NewDataTypeNode(pos *tokenize.Position, datatype *DataType) *Node {
	n := NewNode(NdDataType, pos)
	n.DataTypeField = &DataTypeField{
		DataType: datatype,
	}
	return n
}

func NewImportNode(pos *tokenize.Position, target string) *Node {
	n := NewNode(NdImport, pos)
	n.ImportField = &ImportField{Target: target}
	return n
}

func NewIdentNode(pos *tokenize.Position, ident string) *Node {
	n := NewNode(NdIdent, pos)
	n.IdentField = &IdentField{Ident: ident}
	return n
}

func NewVarDeclNode(pos *tokenize.Position, ident, type_ *Node) *Node {
	n := NewNode(NdVarDecl, pos)
	n.VarDeclField = &VarDeclField{
		Identifier: ident,
		Type:       type_,
	}
	return n
}

func NewAssignNode(pos *tokenize.Position, to, value *Node) *Node {
	n := NewNode(NdAssign, pos)
	n.AssignField = &AssignField{
		To:    to,
		Value: value,
	}
	return n
}

func NewCommentNode(pos *tokenize.Position, comment string) *Node {
	n := NewNode(NdComment, pos)
	n.CommentField = &CommentField{Comment: comment}
	return n
}

func NewFuncDefNode(pos *tokenize.Position, ident, params, returns, body *Node) *Node {
	n := NewNode(NdFuncDef, pos)
	n.FuncDefField = &FuncDefField{
		Identifier: ident,
		Parameters: params,
		Returns:    returns,
		Body:       body,
	}
	return n
}

func NewBlockNode(pos *tokenize.Position, stmts []*Node) *Node {
	n := NewNode(NdBlock, pos)
	n.BlockField = &BlockField{Statements: stmts}
	return n
}

func NewReturnNode(pos *tokenize.Position, value *Node) *Node {
	n := NewNode(NdReturn, pos)
	n.ReturnField = &ReturnField{Value: value}
	return n
}

func NewIfElseNode(pos *tokenize.Position, useElse bool, cond, if_, else_ *Node) *Node {
	n := NewNode(NdIfElse, pos)
	n.IfElseField = &IfElseField{
		UseElse:   useElse,
		Cond:      cond,
		IfBlock:   if_,
		ElseBlock: else_,
	}
	return n
}

func NewForNode(pos *tokenize.Position, init, cond, loop, body *Node) *Node {
	n := NewNode(NdFor, pos)
	n.ForField = &ForField{
		Init: init,
		Cond: cond,
		Loop: loop,
		Body: body,
	}
	return n
}

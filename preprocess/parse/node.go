package parse

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
)

type Node struct {
	Kind NodeKind
	Pos  *tokenize.Position

	//Lhs      *Node
	//Rhs      *Node
	//Values []*Node

	ImportField       *ImportField
	DataTypeField     *DataTypeField
	IdentField        *IdentField
	VarDeclField      *VarDeclField
	AssignField       *AssignField
	CommentField      *CommentField
	FuncDefField      *FuncDefField
	BlockField        *BlockField
	ReturnField       *ReturnField
	IfElseField       *IfElseField
	ForField          *ForField
	ShortVarDeclField *ShortVarDeclField
	BinaryField       *BinaryField
	UnaryField        *UnaryField
	LiteralField      *LiteralField
	CallField         *CallField
	PolynomialField   *PolynomialField
	FuncParam         *FuncParam
}

func (n *Node) String() string {
	var s string
	switch n.Kind {
	case NdImport:
		s = fmt.Sprintf("%v", n.ImportField)
	case NdDataType:
		s = fmt.Sprintf("%v", n.DataTypeField)
	case NdIdent:
		s = fmt.Sprintf("%v", n.IdentField)
	case NdVarDecl:
		s = fmt.Sprintf("%v", n.VarDeclField)
	case NdAssign:
		s = fmt.Sprintf("%v", n.AssignField)
	case NdComment:
		s = fmt.Sprintf("%v", n.CommentField)
	case NdFuncDef:
		s = fmt.Sprintf("%v", n.FuncDefField)
	case NdBlock:
		s = fmt.Sprintf("%v", n.BlockField)
	case NdReturn:
		s = fmt.Sprintf("%v", n.ReturnField)
	case NdIfElse:
		s = fmt.Sprintf("%v", n.IfElseField)
	case NdFor:
		s = fmt.Sprintf("%v", n.ForField)
	case NdShortVarDecl:
		s = fmt.Sprintf("%v", n.ShortVarDeclField)
	case NdAnd, NdOr, NdEq, NdNe, NdLt, NdLe, NdGt, NdGe, NdAdd, NdSub, NdMul, NdDiv, NdMod:
		s = fmt.Sprintf("%v", n.BinaryField)
	case NdNot, NdParenthesis:
		s = fmt.Sprintf("%v", n.UnaryField)
	case NdLiteral:
		s = fmt.Sprintf("%v", n.LiteralField)
	case NdCall:
		s = fmt.Sprintf("%v", n.CallField)
	case NdArgs, NdParams, NdReturnTypes:
		s = fmt.Sprintf("%v", n.PolynomialField)
	case NdParam:
		s = fmt.Sprintf("%v", n.FuncParam)
	}
	return fmt.Sprintf("Node(%d-%d) %s", n.Pos.LineNo, n.Pos.Lat, s)
}

func NewNode(kind NodeKind, pos *tokenize.Position) *Node {
	return &Node{
		Kind: kind,
		Pos:  pos,
	}
}

//func NewNodeWithLR(kind NodeKind, pos *tokenize.Position, lhs, rhs *Node) *Node {
//	return &Node{
//		Kind: kind,
//		Pos:  pos,
//		Lhs:  lhs,
//		Rhs:  rhs,
//	}
//}
//
//func NewNodeWithChildren(kind NodeKind, pos *tokenize.Position, children []*Node) *Node {
//	return &Node{
//		Kind:     kind,
//		Pos:      pos,
//		Values: children,
//	}
//}

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

//func NewReturnNode(pos *tokenize.Position, value *Node) *Node {
//	n := NewNode(NdReturn, pos)
//	n.ReturnField = &ReturnField{Value: value}
//	return n
//}

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

func NewShortVarDeclNode(pos *tokenize.Position, ident, value *Node) *Node {
	n := NewNode(NdShortVarDecl, pos)
	n.ShortVarDeclField = &ShortVarDeclField{
		Identifier: ident,
		Value:      value,
	}
	return n
}

func NewBinaryNode(kind NodeKind, pos *tokenize.Position, lhs, rhs *Node) *Node {
	n := NewNode(kind, pos)
	n.BinaryField = &BinaryField{
		Lhs: lhs,
		Rhs: rhs,
	}
	return n
}

func NewUnaryNode(kind NodeKind, pos *tokenize.Position, value *Node) *Node {
	n := NewNode(kind, pos)
	n.UnaryField = &UnaryField{
		Value: value,
	}
	return n
}

func NewLiteralNode(pos *tokenize.Position, literal *tokenize.Literal) *Node {
	n := NewNode(NdLiteral, pos)
	n.LiteralField = &LiteralField{Literal: literal}
	return n
}

func NewCallNode(pos *tokenize.Position, ident, args *Node) *Node {
	n := NewNode(NdCall, pos)
	n.CallField = &CallField{
		Identifier: ident,
		Args:       args,
	}
	return n
}

func NewPolynomialNode(kind NodeKind, pos *tokenize.Position, values []*Node) *Node {
	n := NewNode(kind, pos)
	n.PolynomialField = &PolynomialField{Values: values}
	return n
}

func NewFuncParamNode(pos *tokenize.Position, ident, typ *Node) *Node {
	n := NewNode(NdParam, pos)
	n.FuncParam = &FuncParam{
		Identifier: ident,
		DataType:   typ,
	}
	return n
}

package parse

type NodeKind int

const (
	_ NodeKind = iota

	NdBlock
	NdReturn
	NdReturnTypes
	NdIf
	NdIfElse
	NdWhile
	NdFor

	NdImport

	NdNot // !
	NdPlus
	NdMinus

	NdAnd // &&
	NdOr  // ||
	NdEq  // ==
	NdNe  // !=
	NdLt  // <
	NdLe  // <=
	NdGt  // >
	NdGe  // >=
	NdAdd // +
	NdSub // -
	NdMul // *
	NdDiv // /
	NdMod // %

	NdFuncDef
	NdVarDecl
	NdShortVarDecl
	NdAssign // =

	NdDataType

	NdLiteral

	NdIdent
	NdCall
	NdFloat
	NdInt
	NdString
	//RawString
	NdList
	NdDict
	NdKV
	NdBool
	NdTrue
	NdFalse
	NdVoid
	NdNil
	NdAny

	NdArgs
	NdParams
	NdParam

	NdAccess
	NdParenthesis

	NdComment
	//White
	//Newline
)

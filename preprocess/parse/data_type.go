package parse

// 型について

// RuntimeDataType 言語として用意された消せない型
type RuntimeDataType int

const (
	Unknown RuntimeDataType = iota
	Int
	Float
	String
	Bool
	Nil // これ型じゃないんだけど、ポインタ実装するまでは型として扱う
)

type DataType struct {
	Ident  string
	Type   RuntimeDataType
	Custom *DataType
}

func GetDataTypeByIdent(ident string) *DataType {
	switch ident {
	case "int":
		return RuntimeInt
	case "float":
		return RuntimeFloat
	case "string":
		return RuntimeString
	case "bool":
		return RuntimeBool
	case "nil":
		return RuntimeNil
	default:
		return &DataType{
			Ident:  ident,
			Type:   Unknown,
			Custom: nil,
		}
	}
}

type DataTypeField struct {
	DataType *DataType
}

// 組み込み型

var RuntimeInt *DataType
var RuntimeFloat *DataType
var RuntimeString *DataType
var RuntimeBool *DataType
var RuntimeNil *DataType

func init() {
	RuntimeInt = &DataType{
		Ident:  "int",
		Type:   Int,
		Custom: nil,
	}

	RuntimeFloat = &DataType{
		Ident:  "float",
		Type:   Float,
		Custom: nil,
	}

	RuntimeString = &DataType{
		Ident:  "string",
		Type:   String,
		Custom: nil,
	}

	RuntimeBool = &DataType{
		Ident:  "bool",
		Type:   Bool,
		Custom: nil,
	}

	RuntimeNil = &DataType{
		Ident:  "nil",
		Type:   Nil,
		Custom: nil,
	}
}

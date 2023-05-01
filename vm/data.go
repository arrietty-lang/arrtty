package vm

import "fmt"

type DataKind int

const (
	KIllegal     = iota
	KLiteral     // string, int, float
	KOffset      // sp, bp, sp+1, ...
	KRegisterTag // R1, R2, ...
	KLabel       // main, f:, ...
	KOpcode
)

func (dk DataKind) String() string {
	switch dk {
	case KIllegal:
		return "KIllegal"
	case KLiteral:
		return "KLiteral"
	case KOffset:
		return "KOffset"
	case KRegisterTag:
		return "KRegisterTag"
	case KLabel:
		return "KLabel"
	case KOpcode:
		return "KOpcode"
	}
	return "KIllegal"
}

type Data struct {
	kind        DataKind
	literal     Literal
	offset      Offset
	registerTag RegisterTag
	label       Label
	opcode      Opcode
}

func (d *Data) GetKind() DataKind {
	return d.kind
}

func (d *Data) String() string {
	var s string
	switch d.kind {
	case KIllegal:
		s = "illegal"
	case KLiteral:
		s = d.literal.String()
	case KOffset:
		s = d.offset.String()
	case KRegisterTag:
		s = d.registerTag.String()
	case KLabel:
		s = d.label.String()
	case KOpcode:
		s = d.opcode.String()
	}

	return fmt.Sprintf("Data{ kind: %s, val: %s }", d.kind.String(), s)
}

func NewLiteralData(literal Literal) *Data {
	return &Data{
		kind:    KLiteral,
		literal: literal,
	}
}

func NewLiteralDataWithRaw[T string | int | float64](v T) *Data {
	return &Data{
		kind:    KLiteral,
		literal: *NewLiteral(v),
	}
}

func NewOffsetData(offset Offset) *Data {
	return &Data{
		kind:   KOffset,
		offset: offset,
	}
}

func NewRegisterTagData(registerTag RegisterTag) *Data {
	return &Data{
		kind:        KRegisterTag,
		registerTag: registerTag,
	}
}

func NewLabelData(label Label) *Data {
	return &Data{
		kind:  KLabel,
		label: label,
	}
}

func NewOpcodeData(opcode Opcode) *Data {
	return &Data{
		kind:   KOpcode,
		opcode: opcode,
	}
}
